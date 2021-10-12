package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// GetMerchantAccountById finds a merchant account by id
func (db *Db) GetMerchantAccountById(ctx context.Context, id uint64, checkAccountActivationStatus bool) (*models.MerchantAccount, error) {
	const operation = "get_business_account_by_id_db_op"
	db.Logger.Info(fmt.Sprintf("get business account by id database operation. id : %d", id))

	tx := db.getMerchantAccountByIdTxFunc(id, checkAccountActivationStatus)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	acc, ok := result.(*models.MerchantAccount)
	if !ok {
		return nil, service_errors.ErrFailedToCastToType
	}

	return acc, nil
}

// getMerchantAccountByIdTxFunc finds the merchant account by id and wraps it in a db tx.
func (db *Db) getMerchantAccountByIdTxFunc(id uint64, checkAccountActivationStatus bool) func(ctx context.Context, tx *gorm.DB) (interface{},
	error) {
	return func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "merchant_account_exists_by_id_tx"
		db.Logger.Info(fmt.Sprintf("get business account by id database tx."))

		if id == 0 {
			return nil, service_errors.ErrInvalidInputArguments
		}

		var account models.MerchantAccountORM
		if err := tx.Where(&models.MerchantAccountORM{Id: id}).First(&account).Error; err != nil {
			return nil, service_errors.ErrAccountDoesNotExist
		}

		if checkAccountActivationStatus {
			if ok := db.AccountActive(&account); !ok {
				return nil, service_errors.ErrAccountExistButInactive
			}
		}

		acct, err := account.ToPB(ctx)
		if err != nil {
			return nil, err
		}

		return &acct, nil
	}
}
