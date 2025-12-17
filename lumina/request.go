package lumina

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/nexula-rg/go-lumina/lumina/internal/config"
	"github.com/nexula-rg/go-lumina/lumina/internal/dto"
)

// Make API LLM request
func (c *Client) MakeRequest(question string) (*Response, error) {
	req := dto.Request{
		Question: question,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("failed to marshal request: %v", err),
		}
	}

	url := fmt.Sprintf("%s%s", config.Url, config.MakeRequest)
	respModel, err := http.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusBadGateway,
			Msg:        fmt.Sprintf("http request error: %v", err),
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return nil, err
	}

	var resp Response
	if err := json.NewDecoder(respModel.Body).Decode(&resp); err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("error of decoding response: %v", err),
		}
	}

	return &resp, nil
}

// Make API LLM request with ability to cancel request by context
func (c *Client) MakeRequestCtx(ctx context.Context, question string) (*Response, error) {
	reqBody := dto.Request{
		Question: question,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("failed to marshal request: %v", err),
		}
	}

	url := fmt.Sprintf("%s%s", config.Url, config.MakeRequest)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("failed to create request: %v", err),
		}
	}
	req.Header.Set("Content-Type", "application/json")

	respModel, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Msg:        fmt.Sprintf("request canceled by context: %v", err),
			}
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &ErrorResponse{
				StatusCode: http.StatusGatewayTimeout,
				Msg:        fmt.Sprintf("request timeout: %v", err),
			}
		}

		return nil, &ErrorResponse{
			StatusCode: http.StatusBadGateway,
			Msg:        fmt.Sprintf("HTTP request error: %v", err),
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return nil, err
	}

	var resp Response
	if err := json.NewDecoder(respModel.Body).Decode(&resp); err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("decode error: %v", err),
		}
	}

	return &resp, nil
}

// Make async API LLM request
func (c *Client) MakeAsyncRequest(answer_uuid string, question string) (*AsyncResponse, error) {
	req := dto.AsyncRequest{
		AnswerUUID: answer_uuid,
		Question:   question,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("failed to marshal request: %v", err),
		}
	}

	url := fmt.Sprintf("%s%s", config.Url, config.MakeAsyncRequest)
	respModel, err := http.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusBadGateway,
			Msg:        fmt.Sprintf("HTTP request error: %v", err),
		}
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return nil, err
	}

	var resp AsyncResponse
	if err := json.NewDecoder(respModel.Body).Decode(&resp); err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("decode response error: %v", err),
		}
	}

	return &resp, nil
}

// Check health condition of API LLM
func (c *Client) CheckHealth() (*HealthResponse, error) {
	url := fmt.Sprintf("%s%s", config.Url, config.CheckHealth)
	respModel, err := http.Get(url)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusBadGateway,
			Msg:        fmt.Sprintf("cannot reach external service: %v", err),
		}
	}
	defer respModel.Body.Close()

	switch respModel.StatusCode {
	case http.StatusOK, http.StatusServiceUnavailable:
	default:
		if err := checkStatusCode(respModel); err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(respModel.Body)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("error reading response body: %v", err),
		}
	}

	var resp HealthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("error unmarshaling json: %v", err),
		}
	}

	return &resp, nil
}

// Ping API LLM
func (c *Client) Ping() error {
	url := fmt.Sprintf("%s%s", config.Url, config.CheckHealth)
	respModel, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot reach external service: %v", err)
	}
	defer respModel.Body.Close()

	if err := checkStatusCode(respModel); err != nil {
		return err
	}

	return nil
}

// checkStatusCode is the function for filter responses from API dto by status codes
func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusNotFound {
			return &ErrorResponse{
				StatusCode: resp.StatusCode,
				Msg:        "http request error: endpoint of external service not founded",
			}
		}

		var respErr ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&respErr); err != nil {
			return &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Msg:        fmt.Sprintf("error of decoding response: %v", err),
			}
		}

		return &ErrorResponse{
			StatusCode: respErr.StatusCode,
			Msg:        fmt.Sprintf("http request error: %v", respErr.Error()),
		}
	}

	return nil
}
