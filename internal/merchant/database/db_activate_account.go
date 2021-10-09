package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// ActivateAccount activates a business account and saves it to the database as long as it exists
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) ActivateAccount(ctx context.Context, id uint64) (bool, error) {
	const operationType = "activate_business_account_db_op"
	db.Logger.Info(fmt.Sprintf("active business account database operation. id: %d", id))

	tx := db.activateMerchantAccountTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return false, err
	}

	opStatus, ok := result.(bool)
	if !ok {
		return false, service_errors.ErrFailedToCastToType
	}

	return opStatus, nil
}

// activateMerchantAccountTxFunc wraps the update operation in a database tx.
func (db *Db) activateMerchantAccountTxFunc(id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operationType = "activate_business_account_db_tx"
		db.Logger.Info("starting transaction")

		if id == 0 {
			return false, service_errors.ErrInvalidInputArguments
		}

		account, err := db.GetMerchantAccountById(ctx, id, false)
		if err != nil {
			return false, err
		}

		if account.IsActive {
			return true, nil
		}

		if err := db.Conn.Engine.Model(&models.MerchantAccountORM{}).Where("id", account.Id).Update("is_active", "true").Error; err != nil {
			return false, err
		}

		return true, nil
	}
	return tx
}
