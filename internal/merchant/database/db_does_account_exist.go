package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// CheckAccountExistenceStatus checks if a merchant account exists solely off its Id
func (db *Db) CheckAccountExistenceStatus(ctx context.Context, id uint64) (bool, error) {
	const operation = "does_business_account_exist_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("get business account existense status database operation. id : %d", id))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.doesMerchantAccountExistTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return true, err
	}

	status, ok := result.(bool)
	if !ok {
		return true, service_errors.ErrFailedToCastToType
	}

	return status, nil
}

// doesMerchantAccountExistTxFunc returns a database transaction wrapping the underlying db logic
func (db *Db) doesMerchantAccountExistTxFunc(id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "does_business_account_exist_db_tx"
		span := db.TracingEngine.CreateChildSpan(ctx, operation)
		defer span.Finish()

		if id == 0 {
			return false, service_errors.ErrInvalidInputArguments
		}

		if ok, err := db.FindMerchantAccountById(ctx, id); !ok && err != nil {
			return false, service_errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
	return tx
}
