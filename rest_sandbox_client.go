package sdk

type SandboxRestClient struct {
	*RestClient
}

func NewSandboxRestClient(token string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, RestApiURL+"/sandbox")}
}

func NewSandboxRestClientCustom(token, apiURL string) *SandboxRestClient {
	return &SandboxRestClient{NewRestClientCustom(token, apiURL)}
}

func (c *SandboxRestClient) Register() error {
	path := c.apiURL + "/sandbox/register"

	return c.postJSONThrow(path, nil)
}

func (c *SandboxRestClient) Clear() error {
	path := c.apiURL + "/sandbox/clear"

	return c.postJSONThrow(path, nil)
}

func (c *SandboxRestClient) SetCurrencyBalance(currency Currency, balance float64) error {
	path := c.apiURL + "/sandbox/currencies/balance"

	payload := struct {
		Currency Currency `json:"currency"`
		Balance  float64  `json:"balance"`
	}{Currency: currency, Balance: balance}

	return c.postJSONThrow(path, payload)
}

func (c *SandboxRestClient) SetPositionsBalance(figi string, balance float64) error {
	path := c.apiURL + "/sandbox/positions/balance"

	payload := struct {
		FIGI    string  `json:"figi"`
		Balance float64 `json:"balance"`
	}{FIGI: figi, Balance: balance}

	return c.postJSONThrow(path, payload)
}
