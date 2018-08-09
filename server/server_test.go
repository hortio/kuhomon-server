package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kumekay/kuhomon-server/device_builder/builder"
	"github.com/kumekay/kuhomon-server/server/db"
	"github.com/kumekay/kuhomon-server/server/model"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	database := db.SetupDB()
	defer database.Close()

	server := NewServer(database)
	router := server.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	type response struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	}

	res := response{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ok", res.Status)
	assert.Equal(t, "Kuhomon HTTP JSON API", res.Description)
}

func TestGetMeasurements(t *testing.T) {
	database := db.SetupDB()
	defer database.Close()

	// Prepare data
	// Device
	device, tokens := builder.BuildDevice(database)

	// Measurements
	m := model.Measurement{
		Pressure:    100000,
		CO2:         500,
		Humidity:    40,
		Temperature: 20,
		At:          time.Now(),
		DeviceID:    device.ID}
	database.Create(&m)

	m = model.Measurement{
		Pressure:    120000,
		CO2:         300,
		Humidity:    20,
		Temperature: 30,
		At:          time.Now(),
		DeviceID:    device.ID}
	database.Create(&m)

	server := NewServer(database)
	router := server.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/measurements/"+device.ID.String(), nil)
	req.Header.Set("Device-Read-Token", tokens.ReadToken)
	router.ServeHTTP(w, req)

	type response struct {
		Measurements []model.AggregatedMeasurement `json:"measurements"`
	}
	res := response{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, len(res.Measurements))
	assert.Equal(t, 400, res.Measurements[0].AvgCO2)
	assert.Equal(t, float32(30), res.Measurements[0].AvgHumidity)
	assert.Equal(t, float32(25), res.Measurements[0].AvgTemperature)
}

func TestPostMeasurements(t *testing.T) {
	database := db.SetupDB()
	defer database.Close()

	// Prepare data
	// Device
	device, tokens := builder.BuildDevice(database)

	m := model.Measurement{
		Pressure:    120000,
		CO2:         300,
		Humidity:    20,
		Temperature: 30}

	server := NewServer(database)
	router := server.setupRouter()

	w := httptest.NewRecorder()
	mJSON, _ := json.Marshal(m)
	req, _ := http.NewRequest("POST", "/measurements/"+device.ID.String(), bytes.NewBuffer(mJSON))
	req.Header.Set("Device-Write-Token", tokens.WriteToken)
	router.ServeHTTP(w, req)

	type response struct {
		At string `json:"at"`
	}
	assert.Equal(t, 201, w.Code)

	res := response{}
	json.Unmarshal(w.Body.Bytes(), &res)
	_, err := time.Parse(
		time.RFC3339,
		res.At)
	assert.Nil(t, err)
}
