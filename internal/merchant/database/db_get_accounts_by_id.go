package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"gorm.io/gorm"
)

// GetMerchantAccountsById obtains a set of business accounts by specified ids
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) GetMerchantAccountsById(ctx context.Context, ids []uint64) ([]*models.MerchantAccount, error) {
	const operation = "get_business_accounts_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("get business account sdatabase operation."))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.getMerchantAccountsTxFunc(ids)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	accounts := result.([]*models.MerchantAccount)
	return accounts, nil
}

// getMerchantAccountsTxFunc wraps the operation in a database tx and returns it.
func (db *Db) getMerchantAccountsTxFunc(ids []uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operationType = "get_business_accounts_db_tx"
		db.Logger.For(ctx).Info("starting database transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, operationType)
		defer span.Finish()

		var accounts = make([]*models.MerchantAccount, len(ids)+1)
		if err := tx.Where(ids).Where(models.MerchantAccount{IsActive: true}).Find(&accounts).Error; err != nil {
			return nil, err
		}

		return accounts, nil
	}
	return tx
}
