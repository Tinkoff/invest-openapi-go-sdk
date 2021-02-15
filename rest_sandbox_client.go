package sdk

import (
	"context"
	"fmt"
)

// SandboxRestClient rest client for sandbox tinkoff invest.
type SandboxRestClient struct {
	*RestClient
}

// NewSandboxRestClient returns new SandboxRestClient by token.
func NewSandboxRestClient(token string) *SandboxRestClient {
	return &SandboxRestClient{RestClient: NewRestClient(token, WithURL(RestAPIURL+"/sandbox"))}
}

// NewSandboxRestClientCustom returns new custom SandboxRestClient by token and api url.
func NewSandboxRestClientCustom(token, apiURL string) *SandboxRestClient {
	return &SandboxRestClient{RestClient: NewRestClient(token, WithURL(apiURL))}
}

// Register see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/sandbox/post_sandbox_register.
func (c *SandboxRestClient) Register(ctx context.Context, accountType AccountType) (Account, error) {
	var response struct {
		Payload Account `json:"payload"`
	}

	path := c.url + "/sandbox/register"

	payload := struct {
		AccountType AccountType `json:"brokerAccountType"`
	}{AccountType: accountType}

	err := c.provider.Post(ctx, path, c.token, payload, &response)
	if err != nil {
		return Account{}, fmt.Errorf("provider post: %w", err)
	}

	return response.Payload, nil
}

// Clear see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/sandbox/post_sandbox_clear.
func (c *SandboxRestClient) Clear(ctx context.Context, accountID string) error {
	path := c.url + "/sandbox/clear"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	err := c.provider.Post(ctx, path, c.token, nil, nil)
	if err != nil {
		return fmt.Errorf("provider post: %w", err)
	}

	return nil
}

// Remove see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/sandbox/post_sandbox_remove.
func (c *SandboxRestClient) Remove(ctx context.Context, accountID string) error {
	path := c.url + "/sandbox/remove"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	err := c.provider.Post(ctx, path, c.token, nil, nil)
	if err != nil {
		return fmt.Errorf("provider post: %w", err)
	}

	return nil
}

// SetCurrencyBalance see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/sandbox/post_sandbox_currencies_balance.
func (c *SandboxRestClient) SetCurrencyBalance(ctx context.Context, accountID string, currency Currency, balance float64) error {
	path := c.url + "/sandbox/currencies/balance"

	payload := struct {
		Currency  Currency `json:"currency"`
		Balance   float64  `json:"balance"`
		AccountID string   `json:"brokerAccountId,omitempty"`
	}{Currency: currency, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	err := c.provider.Post(ctx, path, c.token, payload, nil)
	if err != nil {
		return fmt.Errorf("provider post: %w", err)
	}

	return nil
}

// SetPositionsBalance see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/sandbox/post_sandbox_positions_balance.
func (c *SandboxRestClient) SetPositionsBalance(ctx context.Context, accountID, figi string, balance float64) error {
	path := c.url + "/sandbox/positions/balance"

	payload := struct {
		FIGI      string  `json:"figi"`
		Balance   float64 `json:"balance"`
		AccountID string  `json:"brokerAccountId,omitempty"`
	}{FIGI: figi, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	err := c.provider.Post(ctx, path, c.token, payload, nil)
	if err != nil {
		return fmt.Errorf("provider post: %w", err)
	}

	return nil
}
