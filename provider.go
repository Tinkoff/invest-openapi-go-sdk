package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var _ Provider = &defaultHTTP{}

type (
	// Provider interface for change provider client.
	Provider interface {
		Get(ctx context.Context, url string, token string, unmarshal interface{}) error
		Post(ctx context.Context, url string, token string, payload, unmarshal interface{}) error
	}

	defaultHTTP struct {
		client *http.Client
	}
)

// Post for implements Provider.
func (c *defaultHTTP) Post(ctx context.Context, url string, token string, payload, unmarshal interface{}) error {
	var body io.ReadWriter
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("json marslal: %w", err)
		}
		body = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("build new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("provider do: %w", err)
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

// Get for implements Provider.
func (c *defaultHTTP) Get(ctx context.Context, url string, token string, unmarshal interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("provider do: %w", err)
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
		return nil, fmt.Errorf("provider client do: %w", err)
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
