package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// DeactivateMerchantAccount deactivates a business account and updates the database
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) DeactivateMerchantAccount(ctx context.Context, id uint64) (bool, error) {
	const operation = "delete_business_account_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("delete business account database operation. id : %d", id))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.deactivateMerchantAccountTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return false, err
	}

	status, ok := result.(*bool)
	if !ok {
		return false, service_errors.ErrFailedToCastToType
	}

	return *status, nil
}

// deactivateMerchantAccountTxFunc wraps the delete operation of a merchant account in a database transaction
func (db *Db) deactivateMerchantAccountTxFunc(id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "delete_business_account_db_tx"
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, operation)
		defer span.Finish()

		if id == 0 {
			return false, service_errors.ErrInvalidInputArguments
		}

		account, err := db.GetMerchantAccountById(ctx, id)
		if err != nil {
			return nil, service_errors.ErrAccountDoesNotExist
		}

		account.IsActive = false
		if err := db.SaveAccountRecord(tx, account); err != nil {
			db.Logger.For(ctx).Error(service_errors.ErrFailedToUpdateAccountActiveStatus, err.Error())
			return false, err
		}

		return true, nil
	}
	return tx
}
