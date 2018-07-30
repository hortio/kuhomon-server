package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Measurement is a type for storing datapoints
type Measurement struct {
	Pressure    int     `json:"p,omitempty"`
	CO2         int     `gorm:"column:co2" json:"co2,omitempty"`
	Humidity    float32 `json:"h,omitempty"`
	Temperature float32 `json:"t,omitempty"`

	DeviceID uuid.UUID `gorm:"type:uuid"`
	Device   Device    `gorm:"ForeignKey:DeviceID"`

	At time.Time `json:"at,omitempty" gorm:"not null"`
}

// AggregatedMeasurement is a struct for Measurement serialization
type AggregatedMeasurement struct {
	AvgPressure    int       `json:"p,omitempty"`
	AvgCO2         int       `gorm:"column:avg_co2" json:"co2,omitempty"`
	AvgHumidity    float32   `json:"h,omitempty"`
	AvgTemperature float32   `json:"t,omitempty"`
	AggregatedAt   time.Time `json:"at"`
}

// Device data
type Device struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ReadTokenHash  string
	WriteTokenHash string
}
