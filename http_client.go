package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var _ HTTP = &defaultHTTP{}

type (
	// HTTP interface for change http client.
	HTTP interface {
		Get(ctx context.Context, url string, header http.Header, unmarshal interface{}) error
		Post(ctx context.Context, url string, header http.Header, body []byte, unmarshal interface{}) error
	}

	defaultHTTP struct {
		client *http.Client
	}
)

// Post for implements HTTP.
func (c *defaultHTTP) Post(ctx context.Context, url string, header http.Header, body []byte, unmarshal interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("build new request: %w", err)
	}
	req.Header = header

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if unmarshal != nil {
		err = json.NewDecoder(resp.Body).Decode(unmarshal)
		if err != nil {
			return fmt.Errorf("decode json: %w", err)
		}
	}

	return nil
}

// Get for implements HTTP.
func (c *defaultHTTP) Get(ctx context.Context, url string, header http.Header, unmarshal interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build new request: %w", err)
	}
	req.Header = header

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(unmarshal)
	if err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

func (c *defaultHTTP) do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		tradingError := TradingError{}
		err := json.NewDecoder(resp.Body).Decode(&tradingError)
		if err != nil {
			return nil, fmt.Errorf("json decode error: %w", err)
		}

		return nil, tradingError
	}

	return resp, nil
}
