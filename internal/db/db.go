package db

import (
	"log"
	"os"

	"github.com/Itaro18/bt_be/internal/models"
	"github.com/Itaro18/bt_be/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	utils.LoadEnv()
	dsn := os.Getenv("DATABASE_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}
	// DB = DB.Debug()
	// Auto-migrate all models
	// Auto-migrate one by one to avoid dependency issues
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("❌ User migration failed:", err)
	}
	if err := DB.AutoMigrate(&models.Customer{}); err != nil {
		log.Fatal("❌ Customer migration failed:", err)
	}
	if err := DB.AutoMigrate(&models.Property{}); err != nil {
		log.Fatal("❌ Property migration failed:", err)
	}
	if err := DB.AutoMigrate(&models.Booking{}); err != nil {
		log.Fatal("❌ Booking migration failed:", err)
	}

	log.Println("✅ Database connected and migrated")
	return nil
}
