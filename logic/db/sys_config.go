package db

import "gorm.io/gorm"

type SysConfig struct {
	gorm.Model

	KeyName  string `gorm:"size:100;not null"`
	KeyValue string `gorm:"size:300;not null"`
}

func FindConfigByKeyName(keyName string) (config *SysConfig, err error) {
	err = db.Where("key_name = ?", keyName).Find(&config).Error
	return
}

func UpdateValueByKeyName(keyName string, newValue string) error {
	ret := db.Table("sys_configs").Where("key_name = ?", keyName).UpdateColumn("key_value", newValue)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}
