package dto

// Request is the structure for parsing request from endpoints: /make_request | /make_async_request
type Request struct {
	Question string `json:"question"`
}

// AsyncRequest is the structure for parsing response from endpoint: /make_async_request
type AsyncRequest struct {
	AnswerUUID string `json:"answer_uuid"`
	Question   string `json:"question"`
}
