package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TxFunc func(ctx context.Context, tx *gorm.DB) (interface{}, error)

type OperationType string

// Database provides an interface which any database tied to this service should implement
type Database interface {
	CreateMerchantAccount(ctx context.Context, account *models.MerchantAccount) (*models.MerchantAccount, error)
	UpdateMerchantAccount(ctx context.Context, id uint64, account *models.MerchantAccount) (*models.MerchantAccount, error)
	DeactivateMerchantAccount(ctx context.Context, id uint64) (bool, error)
	GetMerchantAccountById(ctx context.Context, id uint64, checkAccountActivationStatus bool) (*models.MerchantAccount, error)
	CheckAccountExistenceStatus(ctx context.Context, id uint64) (bool, error)
	ActivateAccount(ctx context.Context, id uint64) (bool, error)
	FindMerchantAccountByStripeAccountId(ctx context.Context, stripeConnectedAccountId string) (*models.MerchantAccount, error)
}

// Db withholds connection to a postgres database as well as a logging handler
type Db struct {
	Conn                   *core_database.DatabaseConn
	Logger                 *zap.Logger
	Saga                   *saga.SagaCoordinator
	MaxConnectionAttempts  int
	MaxRetriesPerOperation int
	RetryTimeOut           time.Duration
	OperationSleepInterval time.Duration
}

var _ Database = (*Db)(nil)

type ConnectionInitializationParams struct {
	ConnectionParams       *helper.DatabaseConnectionParams
	Logger                 *zap.Logger
	MaxConnectionAttempts  int
	MaxRetriesPerOperation int
	RetryTimeOut           time.Duration
	RetrySleepInterval     time.Duration
}

// New creates a database connection and returns the connection object
func New(ctx context.Context, params *ConnectionInitializationParams) (*Db,
	error) {
	// TODO: generate a span for the database connection attempt
	if params == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if params.ConnectionParams == nil || params.Logger == nil {
		return nil, errors.New(fmt.Sprintf("%s - invalid connection params objects or logger", service_errors.ErrInvalidInputArguments))
	}

	logger := params.Logger
	databaseModels := models.DatabaseModels()

	conn, err := helper.ConnectToDatabase(ctx, params.ConnectionParams, params.Logger, databaseModels...)
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to the database")

	return &Db{
		Conn:                   conn,
		Logger:                 logger,
		Saga:                   saga.NewSagaCoordinator(logger),
		MaxConnectionAttempts:  params.MaxConnectionAttempts,
		MaxRetriesPerOperation: params.MaxRetriesPerOperation,
		RetryTimeOut:           params.RetryTimeOut,
		OperationSleepInterval: params.RetrySleepInterval,
	}, nil
}
