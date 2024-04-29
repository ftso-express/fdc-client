package restServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

// type ApiResStatusEnum string

// const (
// 	ApiResStatusOk                ApiResStatusEnum = "OK"
// 	ApiResStatusError             ApiResStatusEnum = "ERROR"
// 	ApiResStatusRequestBodyError  ApiResStatusEnum = "REQUEST_BODY_ERROR"
// 	ApiResStatusValidationError   ApiResStatusEnum = "VALIDATION_ERROR"
// 	ApiResStatusTooManyRequests   ApiResStatusEnum = "TOO_MANY_REQUESTS"
// 	ApiResStatusUnauthorized      ApiResStatusEnum = "UNAUTHORIZED"
// 	ApiResStatusAuthError         ApiResStatusEnum = "AUTH_ERROR"
// 	ApiResStatusUpstreamHttpError ApiResStatusEnum = "UPSTREAM_HTTP_ERROR"
// 	ApiResStatusInvalidRequest    ApiResStatusEnum = "INVALID_REQUEST"
// 	ApiResStatusNotImplemented    ApiResStatusEnum = "NOT_IMPLEMENTED"
// 	ApiResStatusPending           ApiResStatusEnum = "PENDING"
// )

// // All api responses should be wrapped in this struct.
// type ApiResponseWrapper[T any] struct {
// 	Data T `json:"data"`

// 	// Optional details for unexpected error responses.
// 	ErrorDetails string `json:"errorDetails"`

// 	// Simple message to explain client developers the reason for error.
// 	ErrorMessage string `json:"errorMessage"`

// 	// Response status. OK for successful responses.
// 	Status ApiResStatusEnum `json:"status"`

// 	ValidationErrorDetails *ApiValidationErrorDetails `json:"validationErrorDetails"`
// }

// type ApiValidationErrorDetails struct {
// 	ClassName   string            `json:"className"`
// 	FieldErrors map[string]string `json:"fieldErrors"`
// }

// writeResponse writes a value to the response writer as a JSON object
// Returns an error if the value could not be written
func writeResponse(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&value)
	if err != nil {
		http.Error(w, fmt.Sprint("error writing response: %w", err), http.StatusInternalServerError)
	}
}

func writeResponseError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		http.Error(w, fmt.Sprint("error writing response: %w", err), http.StatusInternalServerError)
	}

}

// Equivalent to WriteApiResponse with status ApiResponseStatusOk
func WriteApiResponseOk[T any](w http.ResponseWriter, value T) {
	writeResponse(w, value)
}

// Use for error responses
func WriteApiResponseError(
	w http.ResponseWriter,
	errorMessage string,
	errorDetails string,
) {
	writeResponseError(w, errorMessage)
	fmt.Println(errorMessage)
	fmt.Println(errorDetails)
}

// Decode body from the request into value.
// Any error is written into the response and false is returned.
// (It is enough to just return from the request handler on false value)
func DecodeBody(w http.ResponseWriter, r *http.Request, value interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&value)
	if err != nil {
		WriteApiResponseError(w,
			"error parsing request body", err.Error())
		return false
	}
	err = validate.Struct(value)
	if err != nil {
		WriteApiResponseError(w,
			"error validating request body", err.Error())
		return false
	}
	return true
}

func DecodeQueryParams(w http.ResponseWriter, r *http.Request, value interface{}) bool {
	decoder := schema.NewDecoder()
	err := decoder.Decode(value, r.URL.Query())
	if err != nil {
		WriteApiResponseError(w,
			"error parsing query params", err.Error())
		return false
	}
	err = validate.Struct(value)
	if err != nil {
		WriteApiResponseError(w,
			"error validating query params", err.Error())
		return false
	}
	return true
}

func BadParamsErrorHandler(err error) *ErrorHandler {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &ErrorHandler{
		Handler: func(w http.ResponseWriter) {
			WriteApiResponseError(w, err.Error(), err.Error())
		},
	}
}

func InternalServerErrorHandler(err error) *ErrorHandler {
	return &ErrorHandler{
		Handler: func(w http.ResponseWriter) {
			WriteApiResponseError(w, "internal server error", err.Error())
		},
	}
}
