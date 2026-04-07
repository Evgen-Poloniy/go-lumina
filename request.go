package lumina

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nexula-rg/go-lumina/internal/config"
	"github.com/nexula-rg/go-lumina/internal/entity"
)

// Make API LLM request
func (c *Client) MakeRequest(question string) (string, error) {
	reqBody := entity.Request{
		Question: question,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERIALIZATION_ERROR",
				Message: fmt.Sprintf("failed to marshal request: %v", err),
			},
		}
	}

	url := config.Url + config.MakeRequest
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVER_ERROR",
				Message: fmt.Sprintf("failed to create request: %v", err),
			},
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("API-KEY %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusServiceUnavailable,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVICE_UNAVAILABLE",
				Message: fmt.Sprintf("cannot reach external service: %v", err),
			},
		}
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp); err != nil {
		return "", err
	}

	var respModel entity.Response
	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "DESERIALIZATION_ERROR",
				Message: fmt.Sprintf("decode error: %v", err),
			},
		}
	}

	return respModel.Data.Answer, nil
}

// Make API LLM request with ability to cancel request by context
func (c *Client) MakeRequestCtx(ctx context.Context, question string) (string, error) {
	reqBody := entity.Request{
		Question: question,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERIALIZATION_ERROR",
				Message: fmt.Sprintf("failed to marshal request: %v", err),
			},
		}
	}

	url := config.Url + config.MakeRequest
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVER_ERROR",
				Message: fmt.Sprintf("failed to create request: %v", err),
			},
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("API-KEY %s", c.apiKey))

	respModel, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details: struct {
					Code    string
					Message string
				}{
					Code:    "CANCELED",
					Message: fmt.Sprintf("request canceled by context: %v", err),
				},
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return "", &ErrorResponse{
				StatusCode: http.StatusGatewayTimeout,
				Details: struct {
					Code    string
					Message string
				}{
					Code:    "TIMEOUT",
					Message: fmt.Sprintf("request timeout: %v", err),
				},
			}
		}

		return "", &ErrorResponse{
			StatusCode: http.StatusServiceUnavailable,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVICE_UNAVAILABLE",
				Message: fmt.Sprintf("cannot reach external service: %v", err),
			},
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return "", err
	}

	var resp entity.Response
	if err := json.NewDecoder(respModel.Body).Decode(&resp); err != nil {
		return "", &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "DESERIALIZATION_ERROR",
				Message: fmt.Sprintf("decode error: %v", err),
			},
		}
	}

	return resp.Data.Answer, nil
}

// Ping API LLM
func (c *Client) Ping() error {
	url := config.Url + config.CheckHealth
	respModel, err := http.Get(url)
	if err != nil {
		return &ErrorResponse{
			StatusCode: http.StatusServiceUnavailable,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVICE_UNAVAILABLE",
				Message: fmt.Sprintf("cannot reach external service: %v", err),
			},
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return err
	}

	return nil
}

// PingCtx checks the external LLM API with context and returns structured errors
func (c *Client) PingCtx(ctx context.Context) error {
	url := config.Url + config.CheckHealth

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVER_ERROR",
				Message: fmt.Sprintf("failed to create request: %v", err),
			},
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details: struct {
					Code    string
					Message string
				}{
					Code:    "CANCELED",
					Message: fmt.Sprintf("request canceled by context: %v", err),
				},
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return &ErrorResponse{
				StatusCode: http.StatusGatewayTimeout,
				Details: struct {
					Code    string
					Message string
				}{
					Code:    "TIMEOUT",
					Message: fmt.Sprintf("request timeout: %v", err),
				},
			}
		}

		return &ErrorResponse{
			StatusCode: http.StatusServiceUnavailable,
			Details: struct {
				Code    string
				Message string
			}{
				Code:    "SERVICE_UNAVAILABLE",
				Message: fmt.Sprintf("cannot reach external service: %v", err),
			},
		}
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp); err != nil {
		return err
	}

	return nil
}

// checkStatusCode is the function for filter responses from API entity by status codes
func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		errResp.StatusCode = resp.StatusCode

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details: struct {
					Code    string
					Message string
				}{
					Code:    "DESERIALIZATION_ERROR",
					Message: fmt.Sprintf("error of decoding response: %v", err),
				},
			}
		}

		return &errResp

	}

	return nil
}
