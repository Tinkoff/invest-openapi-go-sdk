package sdk

type SandboxTradingClient struct {
	*TradingClient
}

func NewSandboxTradingClient(token string) *SandboxTradingClient {
	return &SandboxTradingClient{NewTradingClientCustom(token, TradingApiURL+"/sandbox")}
}

func NewSandboxTradingClientCustom(token, apiURL string) *SandboxTradingClient {
	return &SandboxTradingClient{NewTradingClientCustom(token, apiURL)}
}

func (c *SandboxTradingClient) Register() error {
	path := c.apiURL + "/sandbox/register"

	return c.postJSONThrow(path, nil)
}

func (c *SandboxTradingClient) Clear() error {
	path := c.apiURL + "/sandbox/clear"

	return c.postJSONThrow(path, nil)
}

func (c *SandboxTradingClient) SetCurrencyBalance(currency Currency, balance float64) error {
	path := c.apiURL + "/sandbox/currencies/balance"

	payload := struct {
		Currency Currency `json:"currency"`
		Balance  float64  `json:"balance"`
	}{Currency: currency, Balance: balance}

	return c.postJSONThrow(path, payload)
}

func (c *SandboxTradingClient) SetPositionsBalance(figi string, balance float64) error {
	path := c.apiURL + "/sandbox/positions/balance"

	payload := struct {
		FIGI    string  `json:"figi"`
		Balance float64 `json:"balance"`
	}{FIGI: figi, Balance: balance}

	return c.postJSONThrow(path, payload)
}
