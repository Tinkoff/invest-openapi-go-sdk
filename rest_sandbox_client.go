package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type SandboxRestClient struct {
	*RestClient
}

func NewSandboxRestClient(token string) *SandboxRestClient {
	tmp, _ := url.Parse(RestApiURL)
	tmp.Path = path.Join(tmp.Path, "/sandbox")
	return &SandboxRestClient{NewRestClientCustom(token, tmp.String())}
}

func NewSandboxRestClientCustom(token, apiURL string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, apiURL)}
}

func (c *SandboxRestClient) Register(ctx context.Context, accountType AccountType) (Account, error) {
	u := c.getUrlRequest("/sandbox/register", nil)
	payload := struct {
		AccountType AccountType `json:"brokerAccountType"`
	}{AccountType: accountType}

	bb, err := json.Marshal(payload)
	if err != nil {
		return Account{}, errors.Errorf("can't marshal request to %s body=%+v", u.String(), payload)
	}

	req, err := c.newRequest(ctx, http.MethodPost, u.String(), bytes.NewReader(bb))
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
		return Account{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", u.String(), respBody)
	}

	return resp.Payload, nil
}

func (c *SandboxRestClient) Clear(ctx context.Context, accountID string) error {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	u := c.getUrlRequest("/sandbox/clear", query)

	return c.postJSONThrow(ctx, u.String(), nil)
}

func (c *SandboxRestClient) Remove(ctx context.Context, accountID string) error {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	u := c.getUrlRequest("/sandbox/remove", query)

	return c.postJSONThrow(ctx, u.String(), nil)
}

func (c *SandboxRestClient) SetCurrencyBalance(ctx context.Context, accountID string, currency Currency, balance float64) error {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	u := c.getUrlRequest("/sandbox/currencies/balance", query)

	payload := struct {
		Currency  Currency `json:"currency"`
		Balance   float64  `json:"balance"`
		AccountID string   `json:"brokerAccountId,omitempty"`
	}{Currency: currency, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	return c.postJSONThrow(ctx, u.String(), payload)
}

func (c *SandboxRestClient) SetPositionsBalance(ctx context.Context, accountID, figi string, balance float64) error {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	u := c.getUrlRequest("/sandbox/positions/balance", query)

	payload := struct {
		FIGI      string  `json:"figi"`
		Balance   float64 `json:"balance"`
		AccountID string  `json:"brokerAccountId,omitempty"`
	}{FIGI: figi, Balance: balance}

	if accountID != DefaultAccount {
		payload.AccountID = accountID
	}

	return c.postJSONThrow(ctx, u.String(), payload)
}
