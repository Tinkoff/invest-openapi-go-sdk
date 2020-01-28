package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type SandboxRestClient struct {
	*RestClient
}

func NewSandboxRestClient(token string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, RestApiURL+"/sandbox")}
}

func NewSandboxRestClientCustom(token, apiURL string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, apiURL)}
}

func (c *SandboxRestClient) Register(ctx context.Context, accountType AccountType) (Account, error) {
	path := c.apiURL + "/sandbox/register"

	payload := struct {
		AccountType AccountType `json:"brokerAccountType"`
	}{AccountType: accountType}

	bb, err := json.Marshal(payload)
	if err != nil {
		return Account{}, errors.Errorf("can't marshal request to %s body=%+v", path, payload)
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(bb))
	if err != nil {
		return Account{}, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return Account{}, err
	}

	type response struct {
		Payload Account `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return Account{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload, nil
}

func (c *SandboxRestClient) Clear(ctx context.Context, accountID string) error {
	path := c.apiURL + "/sandbox/clear"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	return c.postJSONThrow(ctx, path, nil)
}

func (c *SandboxRestClient) Remove(ctx context.Context, accountID string) error {
	path := c.apiURL + "/sandbox/remove"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	return c.postJSONThrow(ctx, path, nil)
}

func (c *SandboxRestClient) SetCurrencyBalance(ctx context.Context, accountID string, currency Currency, balance float64) error {
	path := c.apiURL + "/sandbox/currencies/balance"

	payload := struct {
		Currency  Currency `json:"currency"`
		Balance   float64  `json:"balance"`
		AccountID string   `json:"brokerAccountId,omitempty"`
	}{Currency: currency, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	return c.postJSONThrow(ctx, path, payload)
}

func (c *SandboxRestClient) SetPositionsBalance(ctx context.Context, accountID, figi string, balance float64) error {
	path := c.apiURL + "/sandbox/positions/balance"

	payload := struct {
		FIGI      string  `json:"figi"`
		Balance   float64 `json:"balance"`
		AccountID string  `json:"brokerAccountId,omitempty"`
	}{FIGI: figi, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	return c.postJSONThrow(ctx, path, payload)
}
