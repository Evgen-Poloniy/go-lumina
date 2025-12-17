package lumina

// AsyncResponse is the structure for parsing response from endpoint: /async_requests/make_request
type AsyncResponse struct {
	AnswerUUID string `json:"answer_uuid"`
	Status     string `json:"status"`
}

// Response is the structure for parsing response from endpoints: /requests/make_request
type Response struct {
	Answer string `json:"answer"`
}

// HealthResponse is the structure for parsing response from endpoint: /health
type HealthResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Components struct {
		Model struct {
			Status  string `json:"status"`
			Details string `json:"details"`
		} `json:"model"`
		Database struct {
			Status  string `json:"status"`
			Details string `json:"details"`
		} `json:"database"`
	} `json:"components"`
	Timestamp string `json:"timestamp"`
}
