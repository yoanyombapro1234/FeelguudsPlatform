package models

func DatabaseModels() []interface{} {
	model := make([]interface{}, 0)
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
