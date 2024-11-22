package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

func GetDataBase() *gorm.DB {
	cnf := GetConfig().DataBase
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.User,
		cnf.Pass,
		cnf.Host,
		cnf.Port,
		cnf.Name, // <-- THIS NAME DB
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		slog.Error("failed make connection to database", err)
	}

	// AUTO MIGRATE
	if err = db.AutoMigrate(
		// INSERT HERE IF U WANT AUTO MIGRATE
		//&domain.User{},
	); err != nil {
		slog.Error("failed auto migrate ", err)
	}

	slog.Info("success migrate")
	return db

}
