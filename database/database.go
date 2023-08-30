package database

import (
	"test_jump/config"

	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

func migrateDb() {
	if err := DB.AutoMigrate(&User{}); err != nil {
		fmt.Printf("Failed to migrate User %v", err)
		os.Exit(1)
	}
	if err := DB.AutoMigrate(&Invoice{}); err != nil {
		fmt.Printf("Failed to migrate Invoice %v", err)
		os.Exit(1)
	}
}

func InitDb() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Europe/Paris",
		config.Config.PgConfig.PgHost,
		config.Config.PgConfig.PgUser,
		config.Config.PgConfig.PgPassword,
		config.Config.PgConfig.PgDbName,
		config.Config.PgConfig.PgPort,
		config.Config.PgConfig.PgSslMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to the DB: %v", err)
		os.Exit(1)
	}

	DB = db

	// Development only. Having this code runs by multiple instances on a DB at the same time will cause issues
	// The solution is to use migration scripts and runs them before deployments
	if config.Config.WsConfig.Mode != "release" {
		migrateDb()
	}
}
