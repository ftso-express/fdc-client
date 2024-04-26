package server

type PDPResponseStatus string

const (
	OK            PDPResponseStatus = "OK"
	NOT_AVAILABLE PDPResponseStatus = "NOT_AVAILABLE"
)

type PDPResponse struct {
	Status         PDPResponseStatus `json:"status"`
	Data           string            `json:"data"`
	AdditionalData string            `json:"additionalData"`
}
