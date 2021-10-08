package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// GetMerchantAccountById finds a merchant account by id
func (db *Db) GetMerchantAccountById(ctx context.Context, id uint64) (*models.MerchantAccountORM, error) {
	const operation = "get_business_account_by_id_db_op"
	db.Logger.Info(fmt.Sprintf("get business account by id database operation. id : %d", id))

	tx := db.getMerchantAccountByIdTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	acc, ok := result.(*models.MerchantAccountORM)
	if !ok {
		return nil, service_errors.ErrFailedToCastToType
	}

	return acc, nil
}

// getMerchantAccountByIdTxFunc gets the merchant account by id operation wrapped in a database transaction.
func (db *Db) getMerchantAccountByIdTxFunc(id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "get_business_account_by_id_db_tx"

		if id == 0 {
			return nil, service_errors.ErrInvalidInputArguments
		}

		var account models.MerchantAccountORM
		if err := tx.Where(&models.MerchantAccountORM{Id: id}).First(&account).Error; err != nil {
			return false, service_errors.ErrAccountDoesNotExist
		}

		if ok := db.AccountActive(&account); !ok {
			return false, service_errors.ErrAccountDoesNotExist
		}

		return &account, nil
	}
	return tx
}
