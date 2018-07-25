package main

import (
	"fmt"
	"os"

	"github.com/hortio/kuhomon-server/model"
	bc "golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func hashToken(token string) (string, error) {
	bytes, err := bc.GenerateFromPassword([]byte(token), 14)
	return string(bytes), err
}

func checkHash(token, hash string) bool {
	err := bc.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

func main() {
	db := setupDB()
	server := NewServer(db)
	router := server.setupRouter()
	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}

func setupDB() *gorm.DB {
	// Connect to DB
	dbURI, dbURIPresent := os.LookupEnv("KUHOMON_DB_URL")

	if !dbURIPresent {
		dbURI = "postgresql://ku@localhost:26257/kuhomon_dev?sslmode=disable"
	}

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&model.Measurement{}, &model.Device{})

	return db
}
