package restServer

import (
	"net/http"

	swagger "github.com/davidebianchi/gswagger"
)

type RouteHandler struct {
	Handler            func(w http.ResponseWriter, r *http.Request)
	SwaggerDefinitions swagger.Definitions
	Method             string // Take from net/http package (MethodGet, MethodPost, etc)
}

// // Decode body from the request into value.
// // Any error is written into the response and false is returned.
// // (It is enough to just return from the request handler on false value)
// func decodeBody(w http.ResponseWriter, r *http.Request, value any) (bool, error) {
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&value)
// 	if err != nil {
// 		return false, err
// 		// WriteApiResponseError(w, ApiResStatusRequestBodyError,
// 		// 	"error parsing request body", err.Error())
// 		// return false
// 	}
// 	// err = validate.Struct(value)
// 	// if err != nil {
// 	// 	WriteApiResponseError(w, ApiResStatusRequestBodyError,
// 	// 		"error validating request body", err.Error())
// 	// 	return false
// 	// }
// 	return true, nil
// }

// // TODO: implement standard controller wrappers
// func routeHandlerWithBody[R interface{}]() {
// 	routeHandler := func(w http.ResponseWriter, r *http.Request) {
// 		var request R
// 		if !DecodeBody(w, r, &request) {
// 			return
// 		}
// 		resp, err := handler(request)
// 		if err != nil {
// 			err.Handler(w)
// 			return
// 		}
// 		responseWriter(w, resp)
// 	}
// 	return routeHandler
// }

// func WrappedRouteHandler[RQ interface{}, RS interface{}](
func WrappedRouteHandler[RS interface{}](
	controllerHandler func(http.ResponseWriter, *http.Request),
	method string,
	responseCode int,
	paramDescriptions map[string]string,
	// requestObject RQ,
	respObject RS,
) RouteHandler {
	// Process parameters if any
	pathParams := make(map[string]swagger.Parameter)
	for name, description := range paramDescriptions {
		pathParams[name] = swagger.Parameter{
			Schema:      &swagger.Schema{Value: ""},
			Description: description,
		}
	}

	wrappedRespObject := ApiResponseWrapper[RS]{Data: respObject, Status: ApiResStatusOk}
	swaggerDefinitions := swagger.Definitions{
		PathParams: pathParams,
		// RequestBody: &swagger.ContentValue{
		// 	Content: swagger.Content{
		// 		"application/json": {Value: requestObject},
		// 	},
		// },
		Responses: map[int]swagger.ContentValue{
			responseCode: {
				Content: swagger.Content{
					"application/json": {Value: wrappedRespObject},
				},
			},
		},
	}
	return RouteHandler{
		Handler:            controllerHandler,
		SwaggerDefinitions: swaggerDefinitions,
		Method:             method,
	}
}
