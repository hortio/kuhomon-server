package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// Postgres extensions for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kumekay/kuhomon-server/server/model"
)

// HashToken creates SHA256 hash
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SetupDB creates DB connection
func SetupDB() *gorm.DB {
	// Connect to DB
	dbURI, dbURIPresent := os.LookupEnv("DB_URL")

	if !dbURIPresent {
		dbURI = "postgresql://ku@localhost:26257/kuhomon_dev?sslmode=disable"
	}

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Measurement{}, &model.Device{})

	// Show detailed logs

	if os.Getenv("LOG_DB_QUERIES") == "true" {
		db.LogMode(true)
	}

	return db
}
