package main

import (
	"net/http"
)

// APIError is a type for all possible errors returned by API
type APIError int

// APIErrorDetails is a struct to hold API error descriptions
type APIErrorDetails struct {
	Code int
	ID   string
	Text string
}

// List of possible errors
const (
	DeviceNotFound APIError = iota + 1
	InvalidToken
	Forbidden
	WrongParameters
	DatabaseError
)

// APIErrorDetailsList is a list of all possible errors with descriptions returned by API
var APIErrorDetailsList = map[APIError]APIErrorDetails{
	DeviceNotFound:  APIErrorDetails{http.StatusNotFound, "DeviceNotFound", "Device Not Found"},
	Forbidden:       APIErrorDetails{http.StatusForbidden, "Forbidden", "Action is forbidden"},
	InvalidToken:    APIErrorDetails{http.StatusBadRequest, "InvalidToken", "Token is not valid JWT"},
	WrongParameters: APIErrorDetails{http.StatusBadRequest, "WrongParameters", "Passed parameters are not expectable"},
	DatabaseError:   APIErrorDetails{http.StatusBadRequest, "DatabaseError", "Cannot perform database request"}}
