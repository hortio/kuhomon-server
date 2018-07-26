package main

import (
	"github.com/kumekay/kuhomon-server/server/db"
)

func main() {
	db := db.SetupDB()
	defer db.Close()
	server := NewServer(db)
	router := server.setupRouter()
	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}
