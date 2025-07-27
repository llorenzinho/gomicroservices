package config

import (
	"fmt"
	"users/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func AutoMigrate(db *gorm.DB) error {
	// Automatically migrate the User model
	return db.AutoMigrate(&models.User{})
}

func CreateDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Check connection
	_, err = db.DB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
