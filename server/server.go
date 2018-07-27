package main

import (
	"net/http"
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

	r.GET("/measurements/:deviceID", s.getMeasurements)
	r.POST("/measurements/:deviceID", s.postMeasurements)

	return r
}

func (s *Server) getMeasurements(c *gin.Context) {
	var device model.Device
	var measurements []model.Measurement

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

	err = s.db.Where(&model.Measurement{DeviceID: deviceID}).Order("at desc").Find(&measurements).Error

	if err != nil {
		renderError(c, DatabaseError)
		return
	}

	measurementsJSON := make([]model.MeasurementJSON, len(measurements))

	for i, measurement := range measurements {
		measurementsJSON[i] = model.MeasurementJSON{
			Pressure:    measurement.Pressure,
			CO2:         measurement.CO2,
			Humidity:    measurement.Humidity,
			Temperature: measurement.Temperature,
			At:          measurement.At}
	}

	c.JSON(http.StatusOK, gin.H{"measurements": measurementsJSON})

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
