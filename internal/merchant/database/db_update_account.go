package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// UpdateMerchantAccount updates a business account and saves it to the database as long as it exists
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) UpdateMerchantAccount(ctx context.Context, id uint64, account *models.MerchantAccount) (
	*models.MerchantAccount, error) {
	const operationType = "update_business_account_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("update business account database operation. id: %d", id))

	ctx, span := db.startRootSpan(ctx, operationType)
	defer span.Finish()

	tx := db.updateMerchantAccountTxFunc(account, id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	createdAccount := result.(*models.MerchantAccount)
	return createdAccount, nil
}

// updateMerchantAccountTxFunc wraps the update operation in a database tx.
func (db *Db) updateMerchantAccountTxFunc(account *models.MerchantAccount, id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operationType = "update_business_account_db_tx"
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, operationType)
		defer span.Finish()

		if err := db.ValidateAccount(ctx, account); err != nil {
			return nil, err
		}

		if ok, err := db.FindMerchantAccountById(ctx, id); !ok && err != nil {
			return nil, service_errors.ErrAccountDoesNotExist
		}

		err := db.SaveAccountRecord(tx, account)
		if err != nil {
			return nil, err
		}

		return &account, nil
	}
	return tx
}
