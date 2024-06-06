package restServer

import (
	"net/http"
)

type ExcludeEndpointStruct struct {
	Path   string
	Method string
}

type AipKeyAuthMiddleware struct {
	KeyName          string
	Keys             []string
	ExcludeEndpoints []ExcludeEndpointStruct
	keyMap           map[string]bool
	excludeMap       map[string]map[string]bool
}

func (keyMiddleware *AipKeyAuthMiddleware) Init() {
	keyMiddleware.keyMap = make(map[string]bool)
	keyMiddleware.excludeMap = make(map[string]map[string]bool)
	for _, key := range keyMiddleware.Keys {
		keyMiddleware.keyMap[key] = true
	}
	for _, endpoint := range keyMiddleware.ExcludeEndpoints {
		if keyMiddleware.excludeMap[endpoint.Path] == nil {
			keyMiddleware.excludeMap[endpoint.Path] = make(map[string]bool)
			keyMiddleware.excludeMap[endpoint.Path][endpoint.Method] = true
		} else {
			keyMiddleware.excludeMap[endpoint.Path][endpoint.Method] = true
		}
	}
}

// Middleware function, which will be called for each request
func (keyMiddleware *AipKeyAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if path and method combination are in exclude list
		if _, found := keyMiddleware.excludeMap[r.URL.Path][r.Method]; found {
			// Endpoint is excluded from api key check
			next.ServeHTTP(w, r)
			return
		}

		// Check the api key
		token := r.Header.Get(keyMiddleware.KeyName)
		if _, found := keyMiddleware.keyMap[token]; found {
			// Api key is valid
			next.ServeHTTP(w, r)
			return
		}

		// Access denied
		errorString := "Forbidden, provide valid " + keyMiddleware.KeyName + " api key"
		http.Error(w, errorString, http.StatusForbidden)
	})
}
