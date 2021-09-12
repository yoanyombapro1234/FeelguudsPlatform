package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"gorm.io/gorm"
)

// UpdateAccountOnboardingStatus updates a business account's onboarding status and saves it to the database as long as it exists
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) UpdateAccountOnboardingStatus(ctx context.Context, id uint64, state models.MerchantAccountState) (
	*models.MerchantAccount, error) {
	const operationType = "update_business_account_onboarding_status_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("update business account onboarding status database operation. id: %d", id))

	ctx, span := db.startRootSpan(ctx, operationType)
	defer span.Finish()

	tx := db.updateMerchantAccountOnboardingStatusTxFunc(id, state)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	createdAccount := result.(*models.MerchantAccount)
	return createdAccount, nil
}

// updateMerchantAccountTxFunc wraps the update operation in a database tx.
func (db *Db) updateMerchantAccountOnboardingStatusTxFunc(id uint64, status models.MerchantAccountState) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operationType = "update_business_account_onboarding_status_db_tx"
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, operationType)
		defer span.Finish()

		acct, err := db.GetMerchantAccountById(ctx, id)
		if err != nil {
			return nil, err
		}

		acct.AccountOnboardingState = status
		if err := db.SaveAccountRecord(tx, acct); err != nil {
			return nil, err
		}

		return acct, nil
	}
	return tx
}
