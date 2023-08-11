package database

import (
	"database/sql"
	"log"
	"net/url"
	"time"

	_ "github.com/iruldev/mini-wallet/src/config"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetDB() *gorm.DB {
	return newDB()
}

func newDB() (DB *gorm.DB) {
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	pass := viper.GetString("DB_PASS")
	name := viper.GetString("DB_NAME")
	tz := viper.GetString("DB_TIMEZONE")

	var dsn = ""
	if user != "" {
		dsn = dsn + user
	}
	if pass != "" {
		dsn = dsn + ":" + pass
	}
	if host != "" {
		dsn = dsn + "@tcp(" + host
		if port != "" {
			dsn = dsn + ":" + port + ")"
		} else {
			dsn = dsn + ")"
		}
	}
	if name != "" {
		dsn = dsn + "/" + name
	}
	dsn = dsn + "?parseTime=true&loc=" + url.QueryEscape(tz)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		Conn:       sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{TablePrefix: viper.GetString("DB_PREFIX")},
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Panic(err)
	}

	DB.Set("gorm:auto_preload", true)
	DB.Session(&gorm.Session{
		AllowGlobalUpdate:    true,
		FullSaveAssociations: false,
	})

	sDB, _ := DB.DB()
	sDB.SetMaxOpenConns(10)
	sDB.SetMaxIdleConns(2)
	sDB.SetConnMaxIdleTime(1 * time.Minute)

	return
}
