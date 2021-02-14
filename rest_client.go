package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Errors.
var (
	ErrDepth    = fmt.Errorf("invalid depth. Should be in interval 0 < x <= %d", MaxOrderbookDepth)
	ErrNotFound = errors.New("not found")
)

const (
	// RestAPIURL contains main api url for tinkoff invest api.
	RestAPIURL = "https://api-invest.tinkoff.ru/openapi"
	// MaxTimeout for http request.
	MaxTimeout = time.Second * 30
)

type (
	// RestClient provide to rest methods from tinkoff invest api.
	RestClient struct {
		http  HTTP
		token string
		url   string
	}

	// BuildOption build options for rest client.
	BuildOption func(*RestClient)
)

// WithClient build rest client by custom http client.
func WithClient(http HTTP) BuildOption {
	return func(client *RestClient) {
		client.http = http
	}
}

// WithURL build rest client by custom api url.
func WithURL(url string) BuildOption {
	return func(client *RestClient) {
		client.url = url
	}
}

// NewRestClient build rest client by option.
func NewRestClient(token string, options ...BuildOption) *RestClient {
	client := &RestClient{
		http: &defaultHTTP{
			client: &http.Client{
				Transport: http.DefaultTransport,
				Timeout:   MaxTimeout,
			},
		},
		token: token,
		url:   RestAPIURL,
	}

	for i := range options {
		options[i](client)
	}

	return client
}

// NewRestClientCustom for backward compatibility only.
// Deprecated: have to use NewRestClient by options.
func NewRestClientCustom(token, apiURL string) *RestClient {
	return NewRestClient(token, WithURL(apiURL))
}

// InstrumentByFIGI see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_search_by_figi.
func (c *RestClient) InstrumentByFIGI(ctx context.Context, figi string) (Instrument, error) {
	var response struct {
		Payload Instrument `json:"payload"`
	}

	path := c.url + "/market/search/by-figi?figi=" + figi
	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return Instrument{}, fmt.Errorf("http get: %w", err)
	}

	return response.Payload, nil
}

// InstrumentByTicker see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_search_by_ticker.
func (c *RestClient) InstrumentByTicker(ctx context.Context, ticker string) ([]Instrument, error) {
	var response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	path := c.url + "/market/search/by-ticker?ticker=" + ticker

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Instruments, nil
}

// Currencies see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_currencies.
func (c *RestClient) Currencies(ctx context.Context) ([]Instrument, error) {
	var response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	path := c.url + "/market/currencies"

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Instruments, nil
}

// ETFs see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_etfs.
func (c *RestClient) ETFs(ctx context.Context) ([]Instrument, error) {
	var response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	path := c.url + "/market/etfs"

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Instruments, nil
}

// Bonds see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_bonds.
func (c *RestClient) Bonds(ctx context.Context) ([]Instrument, error) {
	var response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	path := c.url + "/market/bonds"

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Instruments, nil
}

// Stocks see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_stocks.
func (c *RestClient) Stocks(ctx context.Context) ([]Instrument, error) {
	var response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	path := c.url + "/market/stocks"

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Instruments, nil
}

// Operations see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/operations/get_operations.
func (c *RestClient) Operations(ctx context.Context, accountID string, from, to time.Time, figi string) ([]Operation, error) {
	var response struct {
		Payload struct {
			Operations []Operation `json:"operations"`
		} `json:"payload"`
	}

	q := url.Values{
		"from": []string{from.Format(time.RFC3339)},
		"to":   []string{to.Format(time.RFC3339)},
	}
	if figi != "" {
		q.Set("figi", figi)
	}
	if accountID != DefaultAccount {
		q.Set("brokerAccountId", accountID)
	}

	path := c.url + "/operations?" + q.Encode()

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Operations, nil
}

// Portfolio contains calls two method for get portfolio info.
func (c *RestClient) Portfolio(ctx context.Context, accountID string) (Portfolio, error) {
	positions, err := c.PositionsPortfolio(ctx, accountID)
	if err != nil {
		return Portfolio{}, err
	}

	currencies, err := c.CurrenciesPortfolio(ctx, accountID)
	if err != nil {
		return Portfolio{}, err
	}

	return Portfolio{
		Currencies: currencies,
		Positions:  positions,
	}, nil
}

// PositionsPortfolio see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/portfolio/get_portfolio.
func (c *RestClient) PositionsPortfolio(ctx context.Context, accountID string) ([]PositionBalance, error) {
	var response struct {
		Payload struct {
			Positions []PositionBalance `json:"positions"`
		} `json:"payload"`
	}

	path := c.url + "/portfolio"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Positions, nil
}

// CurrenciesPortfolio see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/portfolio/get_portfolio_currencies.
func (c *RestClient) CurrenciesPortfolio(ctx context.Context, accountID string) ([]CurrencyBalance, error) {
	var response struct {
		Payload struct {
			Currencies []CurrencyBalance `json:"currencies"`
		} `json:"payload"`
	}

	path := c.url + "/portfolio/currencies"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Currencies, nil
}

// OrderCancel see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/orders/post_orders_cancel.
func (c *RestClient) OrderCancel(ctx context.Context, accountID, id string) error {
	path := c.url + "/orders/cancel?orderId=" + id
	if accountID != DefaultAccount {
		path += "&brokerAccountId=" + accountID
	}

	err := c.http.Post(ctx, path, c.header(c.token), nil, nil)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}

	return nil
}

// LimitOrder see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/orders/post_orders_limit_order.
func (c *RestClient) LimitOrder(
	ctx context.Context,
	accountID, figi string,
	lots int,
	operation OperationType,
	price float64,
) (PlacedOrder, error) {
	var response struct {
		Payload PlacedOrder `json:"payload"`
	}

	path := c.url + "/orders/limit-order?figi=" + figi

	if accountID != DefaultAccount {
		path += "&brokerAccountId=" + accountID
	}

	payload := struct {
		Lots      int           `json:"lots"`
		Operation OperationType `json:"operation"`
		Price     float64       `json:"price"`
	}{Lots: lots, Operation: operation, Price: price}
	buf, err := json.Marshal(payload)
	if err != nil {
		return PlacedOrder{}, fmt.Errorf("json marshal payload: %w", err)
	}

	err = c.http.Post(ctx, path, c.header(c.token), buf, &response)
	if err != nil {
		return PlacedOrder{}, fmt.Errorf("http post: %w", err)
	}

	return response.Payload, nil
}

// MarketOrder see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/orders/post_orders_market_order.
func (c *RestClient) MarketOrder(ctx context.Context, accountID, figi string, lots int, operation OperationType) (PlacedOrder, error) {
	var response struct {
		Payload PlacedOrder `json:"payload"`
	}

	path := c.url + "/orders/market-order?figi=" + figi

	if accountID != DefaultAccount {
		path += "&brokerAccountId=" + accountID
	}

	payload := struct {
		Lots      int           `json:"lots"`
		Operation OperationType `json:"operation"`
	}{Lots: lots, Operation: operation}

	buf, err := json.Marshal(payload)
	if err != nil {
		return PlacedOrder{}, fmt.Errorf("json marshal payload: %w", err)
	}

	err = c.http.Post(ctx, path, c.header(c.token), buf, &response)
	if err != nil {
		return PlacedOrder{}, fmt.Errorf("http post: %w", err)
	}

	return response.Payload, nil
}

// Orders see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/orders/get_orders.
func (c *RestClient) Orders(ctx context.Context, accountID string) ([]Order, error) {
	var response struct {
		Payload []Order `json:"payload"`
	}

	path := c.url + "/orders"

	if accountID != DefaultAccount {
		path += "?brokerAccountId=" + accountID
	}

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload, nil
}

// Candles see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_candles.
func (c *RestClient) Candles(ctx context.Context, from, to time.Time, interval CandleInterval, figi string) ([]Candle, error) {
	var response struct {
		Payload struct {
			FIGI     string         `json:"figi"`
			Interval CandleInterval `json:"interval"`
			Candles  []Candle       `json:"candles"`
		} `json:"payload"`
	}

	q := url.Values{
		"from":     []string{from.Format(time.RFC3339)},
		"to":       []string{to.Format(time.RFC3339)},
		"interval": []string{string(interval)},
		"figi":     []string{figi},
	}
	path := c.url + "/market/candles?" + q.Encode()

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Candles, nil
}

// Orderbook see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/market/get_market_orderbook.
func (c *RestClient) Orderbook(ctx context.Context, depth int, figi string) (RestOrderBook, error) {
	var response struct {
		Payload RestOrderBook `json:"payload"`
	}

	if depth < 1 || depth > MaxOrderbookDepth {
		return RestOrderBook{}, ErrDepth
	}

	q := url.Values{
		"depth": []string{strconv.Itoa(depth)},
		"figi":  []string{figi},
	}
	path := c.url + "/market/orderbook?" + q.Encode()

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return RestOrderBook{}, fmt.Errorf("http get: %w", err)
	}

	return response.Payload, nil
}

// Accounts see docs https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/#/user/get_user_accounts.
func (c *RestClient) Accounts(ctx context.Context) ([]Account, error) {
	var response struct {
		Payload struct {
			Accounts []Account `json:"accounts"`
		} `json:"payload"`
	}

	path := c.url + "/user/accounts"

	err := c.http.Get(ctx, path, c.header(c.token), &response)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	return response.Payload.Accounts, nil
}

func (c *RestClient) header(token string) http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+token)

	return header
}
