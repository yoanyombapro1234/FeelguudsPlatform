package helper

import (
	"context"
	"errors"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DatabaseConnectionParams struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         int
}

// ConnectToDatabase establish and connects to a database instance
func ConnectToDatabase(ctx context.Context, params *DatabaseConnectionParams, log *zap.Logger, models ...interface{}) (*core_database.DatabaseConn, error) {
	var dbConn *core_database.DatabaseConn
	connectionString := configureConnectionString(ctx, params.Host, params.User, params.Password, params.DatabaseName, params.Port)

	if dbConn = core_database.NewDatabaseConn(connectionString, "postgres"); dbConn == nil {
		return nil, errors.New("failed to connect to merchant component database")
	}

	if err := pingDatabase(ctx, dbConn); err != nil {
		return nil, err
	}

	if err := migrateSchemas(ctx, dbConn, log, models...); err != nil {
		return nil, err
	}

	return dbConn, nil
}

// pingDatabase pings a database node via a connection object
func pingDatabase(ctx context.Context, dbConn *core_database.DatabaseConn) error {
	db, err := dbConn.Engine.DB()
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	setConnectionConfigs(ctx, dbConn)
	return nil
}

// configureConnectionString constructs database connection string from a set of params
func configureConnectionString(ctx context.Context, host, user, password, dbname string, port int) string {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return connectionString
}

// configureDatabaseConnection configures a database connection
func setConnectionConfigs(ctx context.Context, dbConn *core_database.DatabaseConn) {
	dbConn.Engine.FullSaveAssociations = true
	dbConn.Engine.SkipDefaultTransaction = false
	dbConn.Engine.PrepareStmt = true
	dbConn.Engine.DisableAutomaticPing = false
	dbConn.Engine = dbConn.Engine.Set("gorm:auto_preload", true)
}

// migrateSchemas creates or updates a given set of model based on a schema
// if it does not exist or migrates the model schemas to the latest version
func migrateSchemas(ctx context.Context, db *core_database.DatabaseConn, log *zap.Logger, models ...interface{}) error {
	var engine *gorm.DB
	if db == nil {
		return service_errors.ErrInvalidInputArguments
	}

	if engine = db.Engine; engine == nil {
		return errors.New("invalid gorm database engine object")
	}

	if len(models) > 0 {
		if err := engine.AutoMigrate(models...); err != nil {
			// TODO: emit metric
			log.Error(err.Error())
			return err
		}
	}

	return nil
}
