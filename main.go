package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-translator-api/docs"
	"github.com/kujilabo/cocotola-translator-api/pkg/config"
	"github.com/kujilabo/cocotola-translator-api/pkg/gateway"
	"github.com/kujilabo/cocotola-translator-api/pkg/handler"
	"github.com/kujilabo/cocotola-translator-api/pkg/usecase"
	libD "github.com/kujilabo/cocotola-translator-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-translator-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/handler/middleware"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	if len(*env) == 0 {
		appEnv := os.Getenv("APP_ENV")
		if len(appEnv) == 0 {
			*env = "local"
		} else {
			*env = appEnv
		}
	}

	logrus.Infof("env: %s", *env)

	go func() {
		sig := <-sigs
		logrus.Info()
		logrus.Info(sig)
		done <- true
	}()

	cfg, db, sqlDB, router, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	authMiddleware := gin.BasicAuth(gin.Accounts{
		cfg.Auth.Username: cfg.Auth.Password,
	})

	azureTranslationClient := gateway.NewAzureTranslationClient(cfg.Azure.SubscriptionKey)
	rf, err := gateway.NewRepositoryFactory(ctx, db, cfg.DB.DriverName)
	if err != nil {
		panic(err)
	}

	adminUsecase, err := usecase.NewAdminUsecase(rf)
	if err != nil {
		panic(err)
	}
	userUsecase, err := usecase.NewUseUsecase(rf, azureTranslationClient)
	if err != nil {
		panic(err)
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("v1")
	{
		v1.Use(authMiddleware)
		{
			admin := v1.Group("admin")
			adminHandler := handler.NewAdminHandler(adminUsecase)
			admin.POST("find", adminHandler.FindTranslations)
			admin.GET("text/:text/pos/:pos", adminHandler.FindTranslationByTextAndPos)
			admin.GET("text/:text", adminHandler.FindTranslationByText)
			admin.PUT("text/:text/pos/:pos", adminHandler.UpdateTranslation)
			admin.DELETE("text/:text/pos/:pos", adminHandler.RemoveTranslation)
			admin.POST("", adminHandler.AddTranslation)
			admin.POST("export", adminHandler.ExportTranslations)
		}
		{
			admin := v1.Group("user")
			userHandler := handler.NewUserHandler(userUsecase)
			admin.POST("dictionary/lookup", userHandler.DictionaryLookup)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "cocotola.com"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}

	gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			done <- true
		}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exiting")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Infof("Server forced to shutdown. err: %v", err)
	}
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *gin.Engine, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init log
	if err := config.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return nil, nil, nil, nil, err
	}

	// init db
	db, sqlDB, err := initDB(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, xerrors.Errorf("failed to InitDB. err: %w", err)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(middleware.NewLogMiddleware())
	router.Use(gin.Recovery())

	if cfg.Debug.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	return cfg, db, sqlDB, router, nil
}

func initDB(cfg *config.DBConfig) (*gorm.DB, *sql.DB, error) {
	switch cfg.DriverName {
	case "sqlite3":
		db, err := libG.OpenSQLite("./" + cfg.SQLite3.File)
		if err != nil {
			return nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, err
		}

		if err := sqlDB.Ping(); err != nil {
			return nil, nil, err
		}

		if err := libG.MigrateSQLiteDB(db); err != nil {
			return nil, nil, err
		}

		return db, sqlDB, nil
	case "mysql":
		db, err := libG.OpenMySQL(cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)
		if err != nil {
			return nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, err
		}

		if err := sqlDB.Ping(); err != nil {
			return nil, nil, err
		}

		if err := libG.MigrateMySQLDB(db); err != nil {
			return nil, nil, xerrors.Errorf("failed to MigrateMySQLDB. err: %w", err)
		}

		return db, sqlDB, nil
	default:
		return nil, nil, libD.ErrInvalidArgument
	}
}
