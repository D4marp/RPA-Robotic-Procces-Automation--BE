package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"rpa-backend/internal/config"
	"rpa-backend/internal/models"
)

func Connect() *gorm.DB {
	dsn := config.MySQLDSN()
	var db *gorm.DB
	var err error

	for i := 1; i <= 15; i++ {
		log.Printf("Connecting to database (attempt %d/15)...", i)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection failed: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("DB connect error after retries: ", err)
	}

	err = db.AutoMigrate(&models.Job{}, &models.Run{})
	if err != nil {
		log.Fatal("Database migration error: ", err)
	}
	log.Println("MySQL database ready:", config.GetEnv("DB_NAME", "rpa"))
	return db
}

func Seed(db *gorm.DB) {
	templates := []models.Job{
		{
			Name:     "Web Scraper Bot (Price Monitor)",
			Script:   "echo \"Starting Web Scraper Bot...\" && sleep 2 && echo \"Fetching product list from target e-commerce store...\" && sleep 3 && echo \"Found 15 items.\" && echo \"Product: Laptop ASUS Zenbook - IDR 18,500,000\" && echo \"Product: iPhone 15 Pro - IDR 21,000,000 (Price Drop -5%)\" && echo \"Exporting to reports/prices_today.csv...\" && sleep 1 && echo \"Web Scraper Bot completed successfully.\"",
			Schedule: "*/30 * * * *",
			Enabled:  true,
		},
		{
			Name:     "Database Backup & Cloud Sync Bot",
			Script:   "echo \"Starting Database Backup Bot...\" && sleep 1 && echo \"Dumping database 'rpa' to /tmp/backup_rpa_$(date +%F).sql...\" && sleep 2 && echo \"Compressing archive to backup_rpa_$(date +%F).tar.gz (4.8MB)...\" && echo \"Uploading to backups/ folder...\" && sleep 3 && echo \"Upload completed. File ID: gdrive_bk_77391b\" && echo \"Database backup completed successfully.\"",
			Schedule: "0 0 * * *",
			Enabled:  true,
		},
		{
			Name:     "Monthly Financial Report Generator",
			Script:   "echo \"Starting Finance Report Generator...\" && sleep 2 && echo \"Retrieving monthly transactions count: 4,520...\" && sleep 2 && echo \"Calculating total revenue: USD 45,210.50\" && echo \"Generating report: reports/financial_report_June_2026.pdf...\" && sleep 3 && echo \"Sending notification email to finance-team@company.com...\" && sleep 1 && echo \"Bot finished successfully.\"",
			Schedule: "0 8 1 * *",
			Enabled:  true,
		},
		{
			Name:     "Social Media Auto Poster",
			Script:   "echo \"Starting Auto Poster Bot...\" && sleep 1 && echo \"Fetching draft posts from content board...\" && sleep 2 && echo \"Posting Tweet: 'Want to scale your business? Automate boring tasks with RPA.'...\" && echo \"Tweet posted successfully. Tweet ID: 199201992\" && sleep 1 && echo \"Finished auto-posting task.\"",
			Schedule: "0 9,15 * * *",
			Enabled:  true,
		},
		{
			Name:     "System Logs Rotation & Cleaner",
			Script:   "echo \"Starting Logs Cleaner Bot...\" && echo \"Scanning logs directory...\" && echo \"Total log size: 1.2GB\" && sleep 2 && echo \"Rotating log files older than 7 days...\" && echo \"Freeing up 840MB of disk space.\" && sleep 1 && echo \"Logs cleanup complete.\"",
			Schedule: "0 1 * * 0",
			Enabled:  true,
		},
	}

	for _, t := range templates {
		var existing models.Job
		err := db.Where("name = ?", t.Name).First(&existing).Error
		if err != nil {
			db.Create(&t)
			log.Println("Seeded template job:", t.Name)
		}
	}
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
