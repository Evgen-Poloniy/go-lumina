package lumina

// Error is the structure for parsing response from endpoints
type ErrorResponse struct {
	StatusCode int
	Details    struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Details.Message
}
