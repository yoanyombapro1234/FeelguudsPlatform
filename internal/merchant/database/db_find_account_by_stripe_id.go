package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// FindMerchantAccountByStripeAccountId finds a merchant account by stripe id
func (db *Db) FindMerchantAccountByStripeAccountId(ctx context.Context, stripeConnectedAccountId string) (*models.MerchantAccount, error) {
	const operation = "find_merchant_account_by_stripe_db_op"
	db.Logger.Info(fmt.Sprintf("get business account by stripe id database operation."))

	tx := db.findMerchantAccountByStripeConnectedAccountIdTxFunc(stripeConnectedAccountId)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	acct, ok := result.(*models.MerchantAccount)
	if !ok {
		return nil, service_errors.ErrFailedToCastToType
	}

	return acct, nil
}

// findMerchantAccountByStripeConnectedAccountIdTxFunc wraps the logic in a db tx and returns it
func (db *Db) findMerchantAccountByStripeConnectedAccountIdTxFunc(stripeAcctId string) func(ctx context.Context,
	tx *gorm.DB) (interface{}, error) {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "find_merchant_account_by_stripe_tx"
		db.Logger.Info(fmt.Sprintf("get business account by stripe id database tx."))

		if stripeAcctId == EMPTY {
			return nil, service_errors.ErrInvalidInputArguments
		}

		var account models.MerchantAccountORM
		if err := tx.Where(&models.MerchantAccountORM{StripeConnectedAccountId: stripeAcctId}).First(&account).Error; err != nil {
			return nil, service_errors.ErrAccountDoesNotExist
		}

		if !account.IsActive {
			return nil, service_errors.ErrAccountExistButInactive
		}

		acct, err := account.ToPB(ctx)
		if err != nil {
			return nil, err
		}

		return &acct, nil
	}
	return tx
}
