package lumina

// ErrorResponse is the structure for parsing response from endpoints
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Msg        string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Msg
}

func (e *ErrorResponse) GetStatusCode() int {
	return e.StatusCode
}
