package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// FindMerchantAccountById finds a merchant account by id
func (db *Db) FindMerchantAccountById(ctx context.Context, id uint64) (bool, error) {
	const operation = "merchant_account_exists_by_id_op"
	db.Logger.Info(fmt.Sprintf("get business account by id database operation."))

	tx := db.findMerchantAccountByIdTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return true, err
	}

	status, ok := result.(*bool)
	if !ok {
		return true, service_errors.ErrFailedToCastToType
	}

	return *status, nil
}

// findMerchantAccountByIdTxFunc finds the merchant account by id and wraps it in a db tx.
func (db *Db) findMerchantAccountByIdTxFunc(id uint64) func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
	return func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "merchant_account_exists_by_id_tx"
		db.Logger.Info(fmt.Sprintf("get business account by id database tx."))

		if id == 0 {
			return false, service_errors.ErrInvalidInputArguments
		}

		var account models.MerchantAccount
		if err := tx.Where(&models.MerchantAccount{Id: id}).First(&account).Error; err != nil {
			return false, service_errors.ErrAccountDoesNotExist
		}

		if ok := db.AccountActive(&account); !ok {
			return false, service_errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
}
