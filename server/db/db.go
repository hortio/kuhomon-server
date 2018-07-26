package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// Postgres extensions for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kumekay/kuhomon-server/server/model"
	"golang.org/x/crypto/bcrypt"
)

// HashToken creates bcrypt hash
func HashToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	return string(bytes), err
}

// CheckHash verifies bcrypt hash
func CheckHash(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// SetupDB creates DB connection
func SetupDB() *gorm.DB {
	// Connect to DB
	dbURI, dbURIPresent := os.LookupEnv("KUHOMON_DB_URL")

	if !dbURIPresent {
		dbURI = "postgresql://ku@localhost:26257/kuhomon_dev?sslmode=disable"
	}

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Measurement{}, &model.Device{})

	return db
}
