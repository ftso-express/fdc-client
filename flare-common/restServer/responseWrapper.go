package restserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

// writeResponse writes a value to the response writer as a JSON object
// Returns an error if the value could not be written
func writeResponse(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&value)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing response: %s", err), http.StatusInternalServerError)
	}
}

func writeResponseError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing response: %s", err), http.StatusInternalServerError)
	}

}

// Equivalent to WriteAPIResponse with status APIResponseStatusOk
func WriteAPIResponseOk[T any](w http.ResponseWriter, value T) {
	writeResponse(w, value)
}

// Use for error responses
func WriteAPIResponseError(
	w http.ResponseWriter,
	errorMessage string,
	errorDetails string,
) {
	writeResponseError(w, errorMessage)
	log.Info(errorMessage)
	log.Info(errorDetails)
}

// Decode body from the request into value.
// Any error is written into the response and false is returned.
// (It is enough to just return from the request handler on false value)
func DecodeBody(w http.ResponseWriter, r *http.Request, value interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&value)
	if err != nil {
		WriteAPIResponseError(w,
			fmt.Sprintf("error parsing request body on endpoint: %v", r.URL.Path), err.Error())
		return false
	}
	err = validate.Struct(value)
	if err != nil {
		WriteAPIResponseError(w,
			fmt.Sprintf("error validating request body on endpoint: %v", r.URL.Path), err.Error())
		return false
	}
	return true
}

func DecodeQueryParams(w http.ResponseWriter, r *http.Request, value interface{}) bool {
	decoder := schema.NewDecoder()
	err := decoder.Decode(value, r.URL.Query())
	if err != nil {
		WriteAPIResponseError(w,
			fmt.Sprintf("error parsing query params on endpoint: %v", r.URL.Path), err.Error())
		return false
	}
	err = validate.Struct(value)
	if err != nil {
		WriteAPIResponseError(w,
			fmt.Sprintf("error validating query params on endpoint: %v", r.URL.Path), err.Error())
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
			WriteAPIResponseError(w, "Error with params", err.Error())
		},
	}
}

func InternalServerErrorHandler(err error) *ErrorHandler {
	return &ErrorHandler{
		Handler: func(w http.ResponseWriter) {
			WriteAPIResponseError(w, "internal server error", err.Error())
		},
	}
}

func ToEarlyErrorHandler(err error) *ErrorHandler {
	return &ErrorHandler{
		Handler: func(w http.ResponseWriter) {
			WriteAPIResponseError(w, "Request to early", err.Error())
		},
	}
}

func NotAvailableErrorHandler(err error) *ErrorHandler {
	return &ErrorHandler{
		Handler: func(w http.ResponseWriter) {
			WriteAPIResponseError(w, "Data not available yet", err.Error())
		},
	}
}
