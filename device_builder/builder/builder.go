package builder

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kumekay/kuhomon-server/server/db"
	"github.com/kumekay/kuhomon-server/server/model"
)

type DeviceTokens struct {
	ReadToken  string
	WriteToken string
}

// BuildDevice inserts new device into DB
func BuildDevice(database *gorm.DB) (model.Device, DeviceTokens) {
	readToken := tokenGenerator()
	writeToken := tokenGenerator()
	readTokenHash := db.HashToken(readToken)
	writeTokenHash := db.HashToken(writeToken)

	device := model.Device{ReadTokenHash: readTokenHash, WriteTokenHash: writeTokenHash}
	if err := database.Create(&device).Error; err != nil {
		log.Fatal("Cannot create new Device")
	}

	return device, DeviceTokens{readToken, writeToken}
}

func tokenGenerator() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("Cannot Generate random token")
	}

	return fmt.Sprintf("%x", b)
}
