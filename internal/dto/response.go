package dto

// Response is the structure for parsing response from endpoint: /chat
type Response struct {
	Data struct {
		Answer string `json:"answer"`
	} `json:"data"`
}
