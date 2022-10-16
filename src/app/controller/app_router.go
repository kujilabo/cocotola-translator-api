package controller

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola-translator-api/src/app/config"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	"github.com/kujilabo/cocotola-translator-api/src/lib/controller/middleware"
)

func NewRouter(adminUsecase usecase.AdminUsecase, userUsecase usecase.UserUsecase, corsConfig cors.Config, appConfig *config.AppConfig, authConfig *config.AuthConfig, debugConfig *config.DebugConfig) *gin.Engine {
	if !debugConfig.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	authMiddleware := gin.BasicAuth(gin.Accounts{
		authConfig.Username: authConfig.Password,
	})

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))
		v1.Use(authMiddleware)
		{
			admin := v1.Group("admin")
			adminHandler := NewAdminHandler(adminUsecase)
			admin.POST("find", adminHandler.FindTranslationsByFirstLetter)
			admin.GET("text/:text/pos/:pos", adminHandler.FindTranslationByTextAndPos)
			admin.GET("text/:text", adminHandler.FindTranslationsByText)
			admin.PUT("text/:text/pos/:pos", adminHandler.UpdateTranslation)
			admin.DELETE("text/:text/pos/:pos", adminHandler.RemoveTranslation)
			admin.POST("", adminHandler.AddTranslation)
			admin.POST("export", adminHandler.ExportTranslations)
		}
		{
			user := v1.Group("user")
			userHandler := NewUserHandler(userUsecase)
			user.GET("dictionary/lookup", userHandler.DictionaryLookup)
		}
	}

	return router
}
