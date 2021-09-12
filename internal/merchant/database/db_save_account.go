package database

import (
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"gorm.io/gorm"
)

// SaveAccountRecord saves a record in the database
func (db *Db) SaveAccountRecord(tx *gorm.DB, account *models.MerchantAccount) error {
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&account).Error; err != nil {
		return err
	}
	return nil
}
