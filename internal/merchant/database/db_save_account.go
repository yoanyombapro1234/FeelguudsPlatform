package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// SaveAccountRecord saves a record in the database
func (db *Db) SaveAccountRecord(ctx context.Context, acct *models.MerchantAccount) error {
	const operation = "save_merchant_account_db_op"
	db.Logger.Info(fmt.Sprintf("save merchant account database operation."))

	if acct == nil {
		return service_errors.ErrInvalidInputArguments
	}

	merchantAcctOrm, err := acct.ToORM(ctx)
	if err != nil {
		return err
	}

	tx := db.saveMerchantAccountTxFunc(&merchantAcctOrm)
	if err := db.Conn.PerformTransaction(ctx, tx); err != nil {
		return err
	}

	return nil
}

// saveMerchantAccountTxFunc wraps the logic in a db tx and returns it
func (db *Db) saveMerchantAccountTxFunc(acct *models.MerchantAccountORM) func(ctx context.Context,
	tx *gorm.DB) error {
	return func(ctx context.Context, tx *gorm.DB) error {
		const operation = "save_merchant_account_tx"
		db.Logger.Info(fmt.Sprintf("save merchant account database tx."))

		if acct == nil {
			return service_errors.ErrInvalidInputArguments
		}

		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&acct).Error; err != nil {
			return err
		}

		return nil
	}
}
