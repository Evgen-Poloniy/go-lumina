package dto

// Request is the structure for parsing request from endpoint: /chat
type Request struct {
	Question string `json:"question"`
}
