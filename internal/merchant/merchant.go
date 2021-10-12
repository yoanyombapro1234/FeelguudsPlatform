package merchant

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/authentication_handler"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/stripe"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateAccountHandler(w http.ResponseWriter, r *http.Request)
	CreateAccountRefreshUrlHandler(w http.ResponseWriter, r *http.Request)
	CreateAccountReturnUrlHandler(w http.ResponseWriter, r *http.Request)
	DeactivateMerchantAccountHandler(w http.ResponseWriter, r *http.Request)
	GetMerchantAccountHandler(w http.ResponseWriter, r *http.Request)
	UpdateMerchantAccountHandler(w http.ResponseWriter, r *http.Request)
	ReactivateMerchantAccountHandler(w http.ResponseWriter, r *http.Request)
}

// AccountComponent encompasess the suite of merchant account features
type AccountComponent struct {
	// Represents the logging entity which this component uses
	Logger *zap.Logger
	// Represents the database connection object this entity utilizes for storage purposes
	Db *database.Db
	// Represents the object used to interact with the stripe api
	StripeComponent *stripe.Component
	// Used to perform operations against the authentication service
	AuthenticationComponent *authentication_handler.AuthenticationComponent
	// Duration of any expected http call
	HttpTimeout time.Duration
	// Base Refresh url used as part of stripe onboarding process
	BaseRefreshUrl string
	// Base Return url use as part of stripe onboarding process
	BaseReturnUrl string
	// Coordinates distributed tx as a set of saga (compensating and non-compensating tx)
	SagaCoordinater *saga.SagaCoordinator
}

// DatabaseConnectionMetadataParams encompasses connection specific retries and all other associated parameters
type DatabaseConnectionMetadataParams struct {
	// Max number of connection attempts to perform against the database on initial connection initiation
	MaxDatabaseConnectionAttempts int
	// Max number of retries per failed connection attempt
	MaxRetriesPerConnectionAttempt int
	// Max time for a retry to take
	RetryTimeout time.Duration
	// Max time to wait in between retry attempts
	RetrySleepInterval time.Duration
}

// AccountParams encompasses necessary fields to boostrap the merchant account component
type AccountParams struct {
	// Object enables operations against the authentication service
	AuthenticationComponent *authentication_handler.AuthenticationComponent
	// Parameters necessary to initiate a database connection
	DatabaseConnectionParams *helper.DatabaseConnectionParams
	// Parameters necessary to configure database connection retry logic
	DatabaseConnectionMetadataParams *DatabaseConnectionMetadataParams
	// Logging utility
	Logger *zap.Logger
	// Api key used to interact with stripe
	StripeApiKey *string
	// Refresh url used as part of the stripe onboarding process
	RefreshUrl *string
	// Return url used as part of the stripe onboarding process
	ReturnUrl *string
	// Maximum timeout value for all operations
	HttpTimeout time.Duration
}

// NewMerchantAccountComponent returns a new instance of the merchant account component
func NewMerchantAccountComponent(params *AccountParams) *AccountComponent {
	if params == nil {
		log.Fatal("failed to initialize merchant account component due to invalid input arguments")
	}

	if params.Logger == nil {
		log.Fatal("invalid input argument - log object cannot be nil")
	}

	if params.DatabaseConnectionParams == nil {
		log.Fatal("invalid input arguments - db connection params cannot be nil")
	}

	if params.StripeApiKey == nil {
		log.Fatal("invalid input argument - stripe api key cannot be nil")
	}

	if params.RefreshUrl == nil || params.ReturnUrl == nil {
		log.Fatal("invalid input argument - refresh url or return url cannot be nil")
	}

	if params.HttpTimeout == 0 {
		log.Fatal("invalid input argument - http timeout value must be set")
	}

	dbInstance, err := database.New(context.Background(), &database.ConnectionInitializationParams{
		ConnectionParams:       params.DatabaseConnectionParams,
		Logger:                 params.Logger,
		MaxConnectionAttempts:  params.DatabaseConnectionMetadataParams.MaxDatabaseConnectionAttempts,
		MaxRetriesPerOperation: params.DatabaseConnectionMetadataParams.MaxRetriesPerConnectionAttempt,
		RetryTimeOut:           params.DatabaseConnectionMetadataParams.RetryTimeout,
		RetrySleepInterval:     params.DatabaseConnectionMetadataParams.RetrySleepInterval,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	stripeComponent, err := stripe.NewStripeComponent(*params.StripeApiKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	sagaCoordinator := saga.NewSagaCoordinator(params.Logger)

	return &AccountComponent{
		Logger:                  params.Logger,
		Db:                      dbInstance,
		StripeComponent:         stripeComponent,
		AuthenticationComponent: params.AuthenticationComponent,
		BaseReturnUrl:           *params.RefreshUrl,
		BaseRefreshUrl:          *params.ReturnUrl,
		SagaCoordinater:         sagaCoordinator,
	}
}

var _ ServiceInterface = (*AccountComponent)(nil)
