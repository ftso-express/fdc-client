package restServer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResStatusEnum string

const (
	ApiResStatusOk                ApiResStatusEnum = "OK"
	ApiResStatusError             ApiResStatusEnum = "ERROR"
	ApiResStatusRequestBodyError  ApiResStatusEnum = "REQUEST_BODY_ERROR"
	ApiResStatusValidationError   ApiResStatusEnum = "VALIDATION_ERROR"
	ApiResStatusTooManyRequests   ApiResStatusEnum = "TOO_MANY_REQUESTS"
	ApiResStatusUnauthorized      ApiResStatusEnum = "UNAUTHORIZED"
	ApiResStatusAuthError         ApiResStatusEnum = "AUTH_ERROR"
	ApiResStatusUpstreamHttpError ApiResStatusEnum = "UPSTREAM_HTTP_ERROR"
	ApiResStatusInvalidRequest    ApiResStatusEnum = "INVALID_REQUEST"
	ApiResStatusNotImplemented    ApiResStatusEnum = "NOT_IMPLEMENTED"
	ApiResStatusPending           ApiResStatusEnum = "PENDING"
)

// All api responses should be wrapped in this struct.
type ApiResponseWrapper[T any] struct {
	Data T `json:"data"`

	// Optional details for unexpected error responses.
	ErrorDetails string `json:"errorDetails"`

	// Simple message to explain client developers the reason for error.
	ErrorMessage string `json:"errorMessage"`

	// Response status. OK for successful responses.
	Status ApiResStatusEnum `json:"status"`

	ValidationErrorDetails *ApiValidationErrorDetails `json:"validationErrorDetails"`
}

type ApiValidationErrorDetails struct {
	ClassName   string            `json:"className"`
	FieldErrors map[string]string `json:"fieldErrors"`
}

// writeResponse writes a value to the response writer as a JSON object
// Returns an error if the value could not be written
func writeResponse(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&value)
	if err != nil {
		http.Error(w, fmt.Sprint("error writing reponse: %w", err), http.StatusInternalServerError)
	}
}

// Writes value as data field in ApiResponse
func WriteApiResponse[T any](w http.ResponseWriter, status ApiResStatusEnum, value T) {
	response := ApiResponseWrapper[T]{
		Status: status,
		Data:   value,
	}
	writeResponse(w, response)
}

// Equivalent to WriteApiResponse with status ApiResponseStatusOk
func WriteApiResponseOk[T any](w http.ResponseWriter, value T) {
	WriteApiResponse(w, ApiResStatusOk, value)
}

// Use for error responses
func WriteApiResponseError(
	w http.ResponseWriter,
	status ApiResStatusEnum,
	errorMessage string,
	errorDetails string,
) {
	response := ApiResponseWrapper[any]{
		Status:       status,
		ErrorDetails: errorDetails,
		ErrorMessage: errorMessage,
	}
	writeResponse(w, response)
}
