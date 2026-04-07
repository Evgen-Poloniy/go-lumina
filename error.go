package lumina

// Error is the structure for parsing response from endpoints
type ErrorResponse struct {
	StatusCode int
	Details    struct {
		Code    string
		Message string
	}
}

func (e *ErrorResponse) Error() string {
	return e.Details.Message
}

var (
	ErrCreateRequest = "cannot create request"
	ErrReachExternalService = "cannot reach external service"
	ErrInvalidAPIKEY = "invalid API-KEY"
)
