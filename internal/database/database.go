package database

import (
	"sre-chat-api/internal/config"
	"sre-chat-api/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(cfg *config.Config) error {
	var err error

	dsn := cfg.Database.DSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Migrate runs database migrations if enabled
func Migrate(enabled bool) error {
	if !enabled {
		log.Println("Migrations are disabled (MIGRATION_ENABLED=false)")
		return nil
	}

	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	log.Println("Running database migrations...")
	err := DB.AutoMigrate(
		&models.Group{},
		&models.Message{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedDefaultGroup creates the default "SRE Bootcamp" group if it doesn't exist
func SeedDefaultGroup() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	var group models.Group
	result := DB.Where("name = ?", "SRE Bootcamp").First(&group)

	if result.Error == gorm.ErrRecordNotFound {
		group = models.Group{
			Name:        "SRE Bootcamp",
			Description: "Default group for SRE Bootcamp exercises",
		}
		if err := DB.Create(&group).Error; err != nil {
			return fmt.Errorf("failed to create default group: %w", err)
		}
		log.Println("Default group 'SRE Bootcamp' created successfully")
	} else if result.Error != nil {
		return fmt.Errorf("failed to check for default group: %w", result.Error)
	}

	return nil
}
