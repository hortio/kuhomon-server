package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/kumekay/kuhomon-server/server/db"
	"github.com/kumekay/kuhomon-server/server/model"
)

func main() {
	database := db.SetupDB()
	defer database.Close()

	readToken := tokenGenerator()
	writeToken := tokenGenerator()
	readTokenHash, _ := db.HashToken(readToken)
	writeTokenHash, _ := db.HashToken(writeToken)

	device := model.Device{ReadTokenHash: readTokenHash, WriteTokenHash: writeTokenHash}
	if err := database.Create(&device).Error; err != nil {
		log.Fatal("Cannot create new Device")
	}

	fmt.Println("New device succesfully created and inteserted to DB")
	fmt.Println("Device ID:", device.ID)
	fmt.Println("Read Token:", readToken)
	fmt.Println("Read Token Hash:", readTokenHash)
	fmt.Println("Write Token:", writeToken)
	fmt.Println("Write Token Hash:", writeTokenHash)
}

func tokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
