package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/log"
)

var db *gorm.DB

func init() {
	initDB()
}

func buildDSN() string {
	dbConfig := config.Get().DB
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DbName,
		dbConfig.Port,
		dbConfig.SSLMode,
		dbConfig.TimeZone,
	)
}

func mustSuccess(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB() {
	mainDB, err := gorm.Open(postgres.New(postgres.Config{DSN: buildDSN(), PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db = mainDB
	mainDB.Logger = zapgorm2.New(log.GetLogger())

	if config.Get().Debug.Verbose {
		db = db.Debug()
	}

	migrate := db

	mustSuccess(migrate.AutoMigrate(&SysConfig{}))
	mustSuccess(migrate.AutoMigrate(&Transaction{}))
	mustSuccess(migrate.AutoMigrate(&Partner{}))
}
