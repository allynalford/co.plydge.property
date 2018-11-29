package model

// *****************************************************************************
// GenericError
// *****************************************************************************

// GenericError table contains the information for each user
type GenericError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
