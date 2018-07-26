package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Measurement is a type for storing datapoints
type Measurement struct {
	Pressure    int       `json:"p" binding:"required"`
	CO2         int       `json:"co2" binding:"required"`
	Humidity    float32   `json:"h" binding:"required"`
	Temperature float32   `json:"t" binding:"required"`
	DeviceID    uuid.UUID `gorm:"type:uuid"`
	At          time.Time `json:"-"`
}

// Device data
type Device struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReadTokenHash  string
	WriteTokenHash string
}
