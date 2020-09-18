package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const RestApiURL = "https://api-invest.tinkoff.ru/openapi"

var ErrNotFound = errors.New("Not found")

type RestClient struct {
	httpClient *http.Client
	token      string
	apiURL     string
}

func NewRestClient(token string) *RestClient {
	return NewRestClientCustom(token, RestApiURL)
}

func NewRestClientCustom(token string, apiURL string) *RestClient {
	return &RestClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token:  token,
		apiURL: apiURL,
	}
}

func (c *RestClient) InstrumentByFIGI(ctx context.Context, figi string) (Instrument, error) {
	query := url.Values{
		"figi": []string{figi},
	}

	type response struct {
		Payload Instrument `json:"payload"`
	}
	var resp response

	respBody, err := c.requestApi(ctx, "/market/search/by-figi", query, http.MethodGet, nil)
	if err != nil {
		return Instrument{}, err
	}

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return Instrument{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/market/search/by-figi", respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) InstrumentByTicker(ctx context.Context, ticker string) ([]Instrument, error) {
	query := url.Values{
		"ticker": []string{ticker},
	}
	respBody, err := c.requestApi(ctx, "/market/search/by-ticker", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/market/search/by-ticker", respBody)
	}

	return resp.Payload.Instruments, nil
}

func (c *RestClient) Currencies(ctx context.Context) ([]Instrument, error) {
	pathString := path.Join(c.apiURL, "/market/currencies")

	return c.instruments(ctx, pathString)
}

func (c *RestClient) ETFs(ctx context.Context) ([]Instrument, error) {
	pathString := path.Join(c.apiURL, "/market/etfs")

	return c.instruments(ctx, pathString)
}

func (c *RestClient) Bonds(ctx context.Context) ([]Instrument, error) {
	pathString := path.Join(c.apiURL, "/market/bonds")

	return c.instruments(ctx, pathString)
}

func (c *RestClient) Stocks(ctx context.Context) ([]Instrument, error) {
	pathString := path.Join(c.apiURL, "/market/stocks")

	return c.instruments(ctx, pathString)
}

func (c *RestClient) instruments(ctx context.Context, path string) ([]Instrument, error) {
	respBody, err := c.requestApi(ctx, path, nil, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Instruments []Instrument `json:"instruments"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload.Instruments, nil
}

func (c *RestClient) Operations(ctx context.Context, accountID string, from, to time.Time, figi string) ([]Operation, error) {
	query := url.Values{
		"from": []string{from.Format(time.RFC3339)},
		"to":   []string{to.Format(time.RFC3339)},
	}
	if figi != "" {
		query.Set("figi", figi)
	}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	respBody, err := c.requestApi(ctx, "/operations", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Operations []Operation `json:"operations"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/operations", respBody)
	}

	return resp.Payload.Operations, nil
}

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

func (c *RestClient) PositionsPortfolio(ctx context.Context, accountID string) ([]PositionBalance, error) {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	respBody, err := c.requestApi(ctx, "/portfolio", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Positions []PositionBalance `json:"positions"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/portfolio", respBody)
	}

	return resp.Payload.Positions, nil
}

func (c *RestClient) CurrenciesPortfolio(ctx context.Context, accountID string) ([]CurrencyBalance, error) {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	respBody, err := c.requestApi(ctx, "/portfolio/currencies", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Currencies []CurrencyBalance `json:"currencies"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/portfolio/currencies", respBody)
	}

	return resp.Payload.Currencies, nil
}

func (c *RestClient) OrderCancel(ctx context.Context, accountID, id string) error {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	query.Set("orderId", id)
	u := c.getUrlRequest("/orders/cancel", query)

	return c.postJSONThrow(ctx, u.String(), nil)
}

func (c *RestClient) LimitOrder(ctx context.Context, accountID, figi string, lots int, operation OperationType, price float64) (PlacedOrder, error) {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	query.Set("figi", figi)

	payload := struct {
		Lots      int           `json:"lots"`
		Operation OperationType `json:"operation"`
		Price     float64       `json:"price"`
	}{Lots: lots, Operation: operation, Price: price}

	bb, err := json.Marshal(payload)
	if err != nil {
		return PlacedOrder{}, errors.Errorf("can't marshal request to %s body=%+v", "/orders/limit-order", payload)
	}

	respBody, err := c.requestApi(ctx, "/orders/limit-order", query, http.MethodPost, bytes.NewReader(bb))
	if err != nil {
		return PlacedOrder{}, err
	}

	type response struct {
		Payload PlacedOrder `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return PlacedOrder{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/orders/limit-order", respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) MarketOrder(ctx context.Context, accountID, figi string, lots int, operation OperationType) (PlacedOrder, error) {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	query.Set("figi", figi)
	payload := struct {
		Lots      int           `json:"lots"`
		Operation OperationType `json:"operation"`
	}{Lots: lots, Operation: operation}

	bb, err := json.Marshal(payload)
	if err != nil {
		return PlacedOrder{}, errors.Errorf("can't marshal request to %s body=%+v", "/orders/market-order", payload)
	}

	respBody, err := c.requestApi(ctx, "/orders/market-order", query, http.MethodPost, bytes.NewReader(bb))
	if err != nil {
		return PlacedOrder{}, err
	}

	type response struct {
		Payload PlacedOrder `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return PlacedOrder{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/orders/market-order", respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) Orders(ctx context.Context, accountID string) ([]Order, error) {
	query := url.Values{}
	if accountID != DefaultAccount {
		query.Set("brokerAccountId", accountID)
	}
	respBody, err := c.requestApi(ctx, "/orders", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload []Order `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/orders", respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) Candles(ctx context.Context, from, to time.Time, interval CandleInterval, figi string) ([]Candle, error) {
	query := url.Values{
		"from":     []string{from.Format(time.RFC3339)},
		"to":       []string{to.Format(time.RFC3339)},
		"interval": []string{string(interval)},
		"figi":     []string{figi},
	}
	respBody, err := c.requestApi(ctx, "/market/candles", query, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			FIGI     string         `json:"figi"`
			Interval CandleInterval `json:"interval"`
			Candles  []Candle       `json:"candles"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/market/candles", respBody)
	}

	return resp.Payload.Candles, nil
}

func (c *RestClient) Orderbook(ctx context.Context, depth int, figi string) (RestOrderBook, error) {
	if depth < 1 || depth > MaxOrderbookDepth {
		return RestOrderBook{}, ErrDepth
	}

	query := url.Values{
		"depth": []string{strconv.Itoa(depth)},
		"figi":  []string{figi},
	}

	respBody, err := c.requestApi(ctx, "/market/orderbook", query, http.MethodGet, nil)
	if err != nil {
		return RestOrderBook{}, err
	}

	type response struct {
		Payload RestOrderBook `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return RestOrderBook{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/market/orderbook", respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) Accounts(ctx context.Context) ([]Account, error) {
	respBody, err := c.requestApi(ctx, "/user/accounts", nil, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload struct {
			Accounts []Account `json:"accounts"`
		} `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", "/user/accounts", respBody)
	}

	return resp.Payload.Accounts, nil
}

func (c *RestClient) postJSONThrow(ctx context.Context, url string, body interface{}) error {
	var bb []byte
	var err error

	if body != nil {
		bb, err = json.Marshal(body)
		if err != nil {
			return errors.Errorf("can't marshal request body to %s", url)
		}
	}

	req, err := c.newRequest(ctx, http.MethodPost, url, bytes.NewReader(bb))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *RestClient) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Errorf("can't create http request to %s", url)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	return req.WithContext(ctx), nil
}

func (c *RestClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "can't do request to %s", req.URL.RawPath)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read response body to %s", req.URL.RawPath)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		var tradingError TradingError
		if err := json.Unmarshal(body, &tradingError); err == nil {
			return nil, tradingError
		}
		return nil, errors.Errorf("bad response to %s code=%d, body=%s", req.URL, resp.StatusCode, body)
	}

	return body, nil
}

func (c *RestClient) getUrlRequest(urlPath string, query url.Values) *url.URL {
	result, _ := url.Parse(c.apiURL)
	result.Path = path.Join(result.Path, urlPath)
	result.RawQuery = query.Encode()
	return result
}

func (c *RestClient) requestApi(ctx context.Context, urlPath string, query url.Values, method string, body io.Reader) ([]byte, error) {
	u := c.getUrlRequest(urlPath, query)
	req, err := c.newRequest(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

type TradingError struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Hint       string
	Payload    struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"payload"`
}

func (t TradingError) Error() string {
	return fmt.Sprintf(
		"TrackingID: %s, Status: %s, Message: %s, Code: %s, Hint: %s",
		t.TrackingID, t.Status, t.Payload.Message, t.Payload.Code, t.Hint,
	)
}

func (t TradingError) NotEnoughBalance() bool {
	return t.Payload.Code == "NOT_ENOUGH_BALANCE"
}

func (t TradingError) InvalidTokenSpace() bool {
	return t.Payload.Message == "Invalid token scopes"
}
