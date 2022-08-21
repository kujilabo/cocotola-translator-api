package gateway

import (
	"database/sql"
	"errors"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	liberrors "github.com/kujilabo/cocotola-translator-api/src/lib/errors"
)

func ConvertDuplicatedError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == 1062 {
		return newErr
	}

	var sqlite3Err sqlite3.Error
	if ok := errors.As(err, &sqlite3Err); ok && int(sqlite3Err.ExtendedCode) == 2067 {
		return newErr
	}

	return err
}

func ConvertRelationError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == 1452 {
		return newErr
	}

	return err
}

func migrateDB(db *gorm.DB, driverName string, withInstance func(sqlDB *sql.DB) (database.Driver, error)) error {
	sqlDB, err := db.DB()
	if err != nil {
		return liberrors.Errorf("failed to db.DB in gateway.migrateDB. err: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return liberrors.Errorf("failed to os.Getwd in gateway.migrateDB. err: %w", err)
	}

	dir := wd + "/sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		return liberrors.Errorf("failed to gateway.withInstance in gateway.migrateDB. err: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+dir, driverName, driver)
	if err != nil {
		return liberrors.Errorf("failed to migrate.NewWithDatabaseInstance in gateway.migrateDB. err: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return liberrors.Errorf("failed to m.Up in gateway.migrateDB. err: %w", err)
	}

	return nil
}
