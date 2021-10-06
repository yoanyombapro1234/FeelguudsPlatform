package models


func DatabaseModels() []interface{} {
	model := make([]interface{}, 1)
	model = append(model,
		           MerchantAccountORM{},
		           SettingsORM{},
		           ItemSoldORM{},
		           AddressORM{},
		           OwnerORM{},
		           TagsORM{},
		           SettingsORM{})
	return model
}
