package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB returns a gorm.DB instance with proper configs
func NewPostgresDB() *gorm.DB {
	dsn := os.Getenv("AUTH_DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Custom logger for better SQL visibility in dev
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Log slow queries
			LogLevel:                  logger.Info,   // Change to Warn/Error in prod
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying *sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}

	// Connection pool settings (tune depending on workload & server resources)
	sqlDB.SetMaxIdleConns(10)           // idle connections
	sqlDB.SetMaxOpenConns(100)          // max open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // recycle connections every 1h

	log.Println("Connected to Postgres successfully")
	return db
}
