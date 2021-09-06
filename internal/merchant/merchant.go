package merchant

import (
	"context"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"
)

type MerchantServiceInterface interface {
	// CreateMerchantAccount(ctx context.Context, merchantAccount *MerchantAccount)
	// UpdateMerchantAccount(ctx context.Context, merchantAccount *MerchantAccount)
	DeleteMerchantAccount(ctx context.Context, merchantAccountID uint32)
	GetMerchantAccount(ctx context.Context, merchantAccountID uint32)
	GetMerchantAccounts(ctx context.Context, merchantAccountIDs []uint32)
	StartMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
	StopMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
	FinalizeMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32)
}

type MerchantAccountComponent struct {
	Logger *zap.Logger
	Conn   *core_database.DatabaseConn
}

func NewMerchantAccountComponent(params *helper.DatabaseConnectionParams, log *zap.Logger) *MerchantAccountComponent {
	if log == nil || params == nil {
		log.Fatal("failed to initialize merchant account component due to invalid input arguments")
	}

	conn := helper.ConnectToDatabase(params, log)
	if conn == nil {
		log.Fatal("failed to connect to database")
	}

	return &MerchantAccountComponent{
		Logger: log,
		Conn:   conn,
	}
}

func (m MerchantAccountComponent) DeleteMerchantAccount(ctx context.Context, merchantAccountID uint32) {
	panic("implement me")
}

func (m MerchantAccountComponent) GetMerchantAccount(ctx context.Context, merchantAccountID uint32) {
	panic("implement me")
}

func (m MerchantAccountComponent) GetMerchantAccounts(ctx context.Context, merchantAccountIDs []uint32) {
	panic("implement me")
}

func (m MerchantAccountComponent) StartMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32) {
	panic("implement me")
}

func (m MerchantAccountComponent) StopMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32) {
	panic("implement me")
}

func (m MerchantAccountComponent) FinalizeMerchantAccountOnboarding(ctx context.Context, merchantAccountID uint32) {
	panic("implement me")
}

var _ MerchantServiceInterface = (*MerchantAccountComponent)(nil)
