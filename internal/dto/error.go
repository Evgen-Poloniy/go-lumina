package dto

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse is the structure for parsing errors
type ErrorResponse struct {
	Error Error `json:"error"`
}
