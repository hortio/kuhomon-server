package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kumekay/kuhomon-server/server/db"
	"github.com/kumekay/kuhomon-server/server/model"
)

// Server is an http server that handles REST requests.
type Server struct {
	db *gorm.DB
}

// NewServer creates a new instance of a Server.
func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()

	switch os.Getenv("APP_ENV") {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r.GET("/measurements/:deviceID", s.getMeasurements)
	r.POST("/measurements/:deviceID", s.postMeasurements)

	return r
}

func (s *Server) getMeasurements(c *gin.Context) {
	var device model.Device

	deviceID, err := uuid.FromString(c.Param("deviceID"))

	if err != nil {
		renderErrorWithText(c, WrongParameters, "Device ID is not valid UUID")
		return
	}

	if s.db.Where(&model.Device{ID: deviceID}).First(&device).RecordNotFound() {
		renderError(c, DeviceNotFound)
		return
	}

	if db.HashToken(c.GetHeader("Device-Read-Token")) != device.ReadTokenHash {
		renderError(c, Forbidden)
		return
	}

	var measurements []model.AggregatedMeasurement
	err = s.db.Raw(`
		SELECT 
			avg(pressure) as avg_pressure, 
			avg(co2) as avg_co2,
			avg(humidity) as avg_humidity,
			avg(temperature) as avg_temperature,
			((extract('epoch', at) / 600)::INT * 600)::TIMESTAMP as aggregated_at
		FROM measurements
		WHERE device_id = ? AND at > ?
		GROUP BY (extract('epoch', at) / 600)::INT
		ORDER BY aggregated_at DESC`,
		deviceID, time.Now().AddDate(0, 0, -1)).Find(&measurements).Error

	if err != nil {
		renderError(c, DatabaseError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"measurements": measurements})
}

func (s *Server) postMeasurements(c *gin.Context) {
	var m model.Measurement
	var device model.Device
	deviceID, err := uuid.FromString(c.Param("deviceID"))

	if err != nil {
		renderErrorWithText(c, WrongParameters, "Device ID is not valid UUID")
		return
	}

	if s.db.Where(&model.Device{ID: deviceID}).First(&device).RecordNotFound() {
		renderError(c, DeviceNotFound)
		return
	}

	if db.HashToken(c.GetHeader("Device-Write-Token")) != device.WriteTokenHash {
		renderError(c, Forbidden)
		return
	}

	if err := c.BindJSON(&m); err != nil {
		renderError(c, WrongParameters)
		return
	}

	at := time.Now()
	m.At = at
	m.DeviceID = deviceID

	if err := s.db.Create(&m).Error; err != nil {
		renderError(c, DatabaseError)
	} else {
		c.JSON(http.StatusCreated, gin.H{"at": at})
	}
}

func renderError(c *gin.Context, errorID APIError) {
	errorDetails := APIErrorDetailsList[errorID]

	c.JSON(errorDetails.Code, gin.H{
		"id":   errorDetails.ID,
		"text": errorDetails.Text})
}

func renderErrorWithText(c *gin.Context, errorID APIError, text string) {
	errorDetails := APIErrorDetailsList[errorID]

	c.JSON(errorDetails.Code, gin.H{"error": gin.H{
		"id":   errorDetails.ID,
		"text": text}})
}
