package merchant

import (
	"context"
	"time"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/authentication_handler"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/stripe"
	"go.uber.org/zap"
)

type MerchantServiceInterface interface {
	CreateMerchantAccount(ctx context.Context, merchantAccount *models.MerchantAccount) (*models.MerchantAccount, error)
	UpdateMerchantAccount(ctx context.Context, merchantAccount *models.MerchantAccount)
	DeleteMerchantAccount(ctx context.Context, merchantAccountID uint32)
	GetMerchantAccount(ctx context.Context, merchantAccountID uint32)
}

// MerchantAccountComponent encompasess the suite of merchant account features
type MerchantAccountComponent struct {
	Logger                  *zap.Logger
	Db                      *database.Db
	StripeComponent         *stripe.StripeComponent
	AuthenticationComponent *authentication_handler.AuthenticationComponent
	// Duration of any expected http call
	HttpTimeout time.Duration
	// Base Refresh url used as part of stripe onboarding process
	BaseRefreshUrl string
	// Base Return url use as part of stripe onboarding process
	BaseReturnUrl string
}

func NewMerchantAccountComponent(params *helper.DatabaseConnectionParams, log *zap.Logger, stripeAPiKey string,
	authCmp *authentication_handler.AuthenticationComponent) *MerchantAccountComponent {
	if log == nil || params == nil {
		log.Fatal("failed to initialize merchant account component due to invalid input arguments")
	}

	dbInstance, err := database.New(context.Background(), &database.ConnectionInitializationParams{
		ConnectionParams:       params,
		Logger:                 log,
		MaxConnectionAttempts:  2,
		MaxRetriesPerOperation: 3,
		RetryTimeOut:           100,
		RetrySleepInterval:     10,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	component, err := stripe.NewStripeComponent(stripeAPiKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &MerchantAccountComponent{
		Logger:                  log,
		Db:                      dbInstance,
		StripeComponent:         component,
		AuthenticationComponent: authCmp,
		BaseReturnUrl:           "http://localhost/v1/merchant-account/return-url",
		BaseRefreshUrl:          "http://localhost/v1/merchant-account/refresh-url",
	}
}

var _ MerchantServiceInterface = (*MerchantAccountComponent)(nil)
