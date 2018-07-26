package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

	r.GET("/measurements", s.getMeasurements)
	r.POST("/measurements", s.postMeasurements)

	return r
}

// TODO: this is temporary stub
func (s *Server) getMeasurements(c *gin.Context) {
	device := c.GetHeader("Device-Read-Token")
	value := make([]string, 3)
	ok := true
	if ok {
		c.JSON(http.StatusOK, gin.H{"measurements": value, "device": device})
	} else {
		renderError(c, DeviceNotFound)
	}
}

func (s *Server) postMeasurements(c *gin.Context) {
	var m model.Measurement

	// TODO: Find Device
	err := true

	if err {
		renderError(c, DeviceNotFound)
		return
	}

	if err := c.BindJSON(&m); err != nil {
		renderError(c, WrongParameters)
		return
	}

	at := time.Now()
	m.At = at

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
