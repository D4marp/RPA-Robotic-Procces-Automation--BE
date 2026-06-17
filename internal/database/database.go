package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"rpa-backend/internal/config"
	"rpa-backend/internal/models"
)

func Connect() *gorm.DB {
	dsn := config.MySQLDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	db.AutoMigrate(&models.Job{}, &models.Run{})
	log.Println("MySQL database ready:", config.GetEnv("DB_NAME", "rpa"))
	return db
}

// ConnectTest uses in-memory SQLite for fast isolated unit tests.
func ConnectTest() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Job{}, &models.Run{})
	return db
}
