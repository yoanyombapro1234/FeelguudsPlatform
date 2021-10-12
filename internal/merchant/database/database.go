package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TxFunc func(ctx context.Context, tx *gorm.DB) (interface{}, error)

type OperationType string

// Interface provides an interface which any database tied to this service should implement
type Interface interface {
	// CreateMerchantAccount creates a merchant account record
	CreateMerchantAccount(ctx context.Context, account *models.MerchantAccount) (*models.MerchantAccount, error)

	// UpdateMerchantAccount updates a merchant account record assuming the record exists
	UpdateMerchantAccount(ctx context.Context, id uint64, account *models.MerchantAccount) (*models.MerchantAccount, error)

	// DeactivateMerchantAccount performs a "soft" delete of the merchant account record
	DeactivateMerchantAccount(ctx context.Context, id uint64) (bool, error)

	// GetMerchantAccountById returns a merchant account record if it exists
	GetMerchantAccountById(ctx context.Context, id uint64, checkAccountActivationStatus bool) (*models.MerchantAccount, error)

	// CheckAccountExistenceStatus asserts a given merchant account exists based on provided merchant account ID
	CheckAccountExistenceStatus(ctx context.Context, id uint64) (bool, error)

	// ActivateAccount activates a merchant account record
	ActivateAccount(ctx context.Context, id uint64) (bool, error)

	// FindMerchantAccountByStripeAccountId attempts to obtain a merchant account record based on the merchant's
	// stripe connected account ID
	FindMerchantAccountByStripeAccountId(ctx context.Context, stripeConnectedAccountId string) (*models.MerchantAccount, error)
}

// Db withholds connection to a postgres database as well as a logging handler
type Db struct {
	// Conn serves as the actual database connection object
	Conn *core_database.DatabaseConn
	// Logger is the logging utility used by this object
	Logger *zap.Logger
	// MaxConnectionAttempts outlines the maximum connection attempts
	// to initiate against the database
	MaxConnectionAttempts int
	// MaxRetriesPerOperation defines the maximum retries to attempt per failed database
	// connection attempt
	MaxRetriesPerOperation int
	// RetryTimeOut defines the maximum time until a retry operation is observed as a
	// timed out operation
	RetryTimeOut time.Duration
	// OperationSleepInterval defines the amount of time between retry operations
	// that the system sleeps
	OperationSleepInterval time.Duration
}

var _ Interface = (*Db)(nil)

// ConnectionInitializationParams represents connection initialization parameters for the database
type ConnectionInitializationParams struct {
	// ConnectionParams outlines database connection parameters
	ConnectionParams *helper.DatabaseConnectionParams
	// Logger is the logging utility used by this object
	Logger *zap.Logger
	// MaxConnectionAttempts outlines the maximum connection attempts
	// to initiate against the database
	MaxConnectionAttempts int
	// MaxRetriesPerOperation defines the maximum retries to attempt per failed database
	// connection attempt
	MaxRetriesPerOperation int
	// RetryTimeOut defines the maximum time until a retry operation is observed as a
	// timed out operation
	RetryTimeOut time.Duration
	// RetrySleepInterval defines the amount of time between retry operations
	// that the system sleeps
	RetrySleepInterval time.Duration
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
		MaxConnectionAttempts:  params.MaxConnectionAttempts,
		MaxRetriesPerOperation: params.MaxRetriesPerOperation,
		RetryTimeOut:           params.RetryTimeOut,
		OperationSleepInterval: params.RetrySleepInterval,
	}, nil
}
