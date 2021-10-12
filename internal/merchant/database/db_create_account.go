package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-database"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"gorm.io/gorm"
)

// CreateMerchantAccount creates a business account and saves it to the database
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) CreateMerchantAccount(ctx context.Context, account *models.MerchantAccount) (*models.MerchantAccount, error) {
	const operation = "create_business_account_db_op"
	db.Logger.Info(fmt.Sprintf("create business account database operation."))

	if account == nil {
		return nil, service_errors.ErrInvalidAccount
	}

	tx := db.createAccountTxFunc(account)

	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	createdAccount := result.(*models.MerchantAccount)
	return createdAccount, nil
}

// createAccountTxFunc encloses the account creation step in a database transaction
func (db *Db) createAccountTxFunc(account *models.MerchantAccount) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "create_business_account_db_tx"
		db.Logger.Info("starting transaction")

		if account == nil {
			return nil, service_errors.ErrInvalidAccount
		}

		acctOrm, err := account.ToORM(ctx)
		if err != nil {
			return nil, err
		}

		if err := db.ValidateAccount(ctx, &acctOrm); err != nil {
			return nil, err
		}

		if _, err := db.FindMerchantAccountByEmail(ctx, acctOrm.BusinessEmail); err == nil {
			return nil, service_errors.ErrAccountAlreadyExist
		}

		if err := db.UpdateAccountOnboardStatus(ctx, &acctOrm); err != nil {
			return nil, err
		}

		acctOrm.IsActive = true

		// ensure it is saved
		if err = db.Conn.Engine.Create(&acctOrm).Error; err != nil {
			return nil, err
		}

		newAccount, err := acctOrm.ToPB(ctx)
		if err != nil {
			return nil, err
		}

		return &newAccount, nil
	}
	return tx
}
