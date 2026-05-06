package lumina

// Error is the structure for parsing response from endpoints
type Error struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrCreateRequest        = "cannot create request"
	ErrReachExternalService = "cannot reach external service"
	ErrInvalidAPIKEY        = "invalid API-KEY"
)
