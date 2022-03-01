package engine

import (
	"github.com/cyops-se/dd-proxy/db"
	"github.com/cyops-se/dd-proxy/types"
)

func GetSetting(key string) (types.KeyValuePair, error) {
	var item types.KeyValuePair
	result := db.DB.First(&item, "key = ?", key)
	return item, result.Error
}

func PutSetting(key string, value string) {
	if item, err := GetSetting(key); err == nil {
		item.Value = value
		db.DB.Save(item)
	} else {
		item := types.KeyValuePair{Key: key, Value: value}
		db.DB.Create(&item)
	}
}

func InitSetting(key string, value string, description string) types.KeyValuePair {
	item, err := GetSetting(key)
	if err != nil {
		item := types.KeyValuePair{Key: key, Value: value, Extra: description}
		db.DB.Create(&item)
	}

	return item
}
