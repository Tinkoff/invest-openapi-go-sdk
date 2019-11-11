package sdk

import "context"

type SandboxRestClient struct {
	*RestClient
}

func NewSandboxRestClient(token string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, RestApiURL+"/sandbox")}
}

func NewSandboxRestClientCustom(token, apiURL string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, apiURL)}
}

func (c *SandboxRestClient) Register(ctx context.Context) error {
	path := c.apiURL + "/sandbox/register"

	return c.postJSONThrow(ctx, path, nil)
}

func (c *SandboxRestClient) Clear(ctx context.Context) error {
	path := c.apiURL + "/sandbox/clear"

	return c.postJSONThrow(ctx, path, nil)
}

func (c *SandboxRestClient) SetCurrencyBalance(ctx context.Context, currency Currency, balance float64) error {
	path := c.apiURL + "/sandbox/currencies/balance"

	payload := struct {
		Currency Currency `json:"currency"`
		Balance  float64  `json:"balance"`
	}{Currency: currency, Balance: balance}

	return c.postJSONThrow(ctx, path, payload)
}

func (c *SandboxRestClient) SetPositionsBalance(ctx context.Context, figi string, balance float64) error {
	path := c.apiURL + "/sandbox/positions/balance"

	payload := struct {
		FIGI    string  `json:"figi"`
		Balance float64 `json:"balance"`
	}{FIGI: figi, Balance: balance}

	return c.postJSONThrow(ctx, path, payload)
}
