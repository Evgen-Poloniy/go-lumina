package lumina

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nexula-rg/go-lumina/internal/config"
	"github.com/nexula-rg/go-lumina/internal/dto"
)

// Make API LLM request
func (c *Client) MakeRequest(question string) (string, error) {
	reqBody := dto.Request{
		Question: question,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERIALIZATION_ERROR",
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
		}
	}

	url := config.Url + config.MakeRequest
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERVER_ERROR",
			Message:    fmt.Sprintf("failed to create request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", &Error{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "SERVICE_UNAVAILABLE",
			Message:    fmt.Sprintf("cannot reach external service: %v", err),
		}
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp); err != nil {
		return "", err
	}

	var respModel dto.Response
	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "DESERIALIZATION_ERROR",
			Message:    fmt.Sprintf("decode error: %v", err),
		}
	}

	return respModel.Data.Answer, nil
}

// Make API LLM request with ability to cancel request by context
func (c *Client) MakeRequestCtx(ctx context.Context, question string) (string, error) {
	reqBody := dto.Request{
		Question: question,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERIALIZATION_ERROR",
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
		}
	}

	url := config.Url + config.MakeRequest
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERVER_ERROR",
			Message:    fmt.Sprintf("failed to create request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	respModel, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", &Error{
				StatusCode: http.StatusInternalServerError,
				Code:       "CANCELED",
				Message:    fmt.Sprintf("request canceled by context: %v", err),
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return "", &Error{
				StatusCode: http.StatusGatewayTimeout,
				Code:       "TIMEOUT",
				Message:    fmt.Sprintf("request timeout: %v", err),
			}
		}

		return "", &Error{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "SERVICE_UNAVAILABLE",
			Message:    fmt.Sprintf("cannot reach external service: %v", err),
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return "", err
	}

	var resp dto.Response
	if err := json.NewDecoder(respModel.Body).Decode(&resp); err != nil {
		return "", &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "DESERIALIZATION_ERROR",
			Message:    fmt.Sprintf("decode error: %v", err),
		}
	}

	return resp.Data.Answer, nil
}

// Ping the API LLM and check API-KEY on the API LLM side
func (c *Client) Check(apiKey string) error {
	url := config.Url + config.Check

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.New(ErrCreateRequest)
	}

	req.Header.Add("X-API-Key", apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.New(ErrReachExternalService)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New(ErrInvalidAPIKEY)
	}

	return nil
}

// Ping the API LLM
func (c *Client) Ping() error {
	url := config.Url + config.Health

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERVER_ERROR",
			Message:    fmt.Sprintf("failed to create request: %v", err),
		}
	}

	respModel, err := c.client.Do(req)
	if err != nil {
		return &Error{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "SERVICE_UNAVAILABLE",
			Message:    fmt.Sprintf("cannot reach external service: %v", err),
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
	url := config.Url + config.Health

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "SERVER_ERROR",
			Message:    fmt.Sprintf("failed to create request: %v", err),
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return &Error{
				StatusCode: http.StatusInternalServerError,
				Code:       "CANCELED",
				Message:    fmt.Sprintf("request canceled by context: %v", err),
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return &Error{
				StatusCode: http.StatusGatewayTimeout,
				Code:       "TIMEOUT",
				Message:    fmt.Sprintf("request timeout: %v", err),
			}
		}

		return &Error{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "SERVICE_UNAVAILABLE",
			Message:    fmt.Sprintf("cannot reach external service: %v", err),
		}
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp); err != nil {
		return err
	}

	return nil
}

// checkStatusCode is the function for filter responses from API dto by status codes
func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp dto.ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return &Error{
				StatusCode: http.StatusInternalServerError,
				Code:       "DESERIALIZATION_ERROR",
				Message:    fmt.Sprintf("error decoding response: %v", err),
			}
		}

		return &Error{
			StatusCode: resp.StatusCode,
			Code:       errResp.Error.Code,
			Message:    errResp.Error.Message,
		}

	}

	return nil
}
