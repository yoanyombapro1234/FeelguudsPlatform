package merchant

import (
	"context"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"go.uber.org/zap"
)

type MerchantServiceInterface interface {
	CreateMerchantAccount(ctx context.Context, merchantAccount *models.MerchantAccount)
	UpdateMerchantAccount(ctx context.Context, merchantAccount *models.MerchantAccount)
	DeleteMerchantAccount(ctx context.Context, merchantAccountID uint32)
	GetMerchantAccount(ctx context.Context, merchantAccountID uint32)
	GetMerchantAccounts(ctx context.Context, merchantAccountIDs []uint32)
	StartMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
	StopMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
	FinalizeMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
}

type MerchantAccountComponent struct {
	Logger *zap.Logger
	Db     *database.Db
}

func NewMerchantAccountComponent(params *helper.DatabaseConnectionParams, log *zap.Logger) *MerchantAccountComponent {
	if log == nil || params == nil {
		log.Fatal("failed to initialize merchant account component due to invalid input arguments")
	}

	dbInstance, err := database.New(context.Background(), database.ConnectionInitializationParams{
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

	return &MerchantAccountComponent{
		Logger: log,
		Db:     dbInstance,
	}
}

var _ MerchantServiceInterface = (*MerchantAccountComponent)(nil)
