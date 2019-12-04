package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

func NewRestClientCustom(token, apiURL string) *RestClient {
	return &RestClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		token:  token,
		apiURL: apiURL,
	}
}

func (c *RestClient) SearchInstrumentByFIGI(figi string) (Instrument, error) {
	path := c.apiURL + "/market/search/by-figi?figi=" + figi

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return Instrument{}, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return Instrument{}, err
	}

	type response struct {
		Payload Instrument `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return Instrument{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) SearchInstrumentByTicker(ticker string) ([]Instrument, error) {
	path := c.apiURL + "/market/search/by-ticker?ticker=" + ticker

	return c.instruments(path)
}

func (c *RestClient) Currencies() ([]Instrument, error) {
	path := c.apiURL + "/market/currencies"

	return c.instruments(path)
}

func (c *RestClient) ETFs() ([]Instrument, error) {
	path := c.apiURL + "/market/etfs"

	return c.instruments(path)
}

func (c *RestClient) Bonds() ([]Instrument, error) {
	path := c.apiURL + "/market/bonds"

	return c.instruments(path)
}

func (c *RestClient) Stocks() ([]Instrument, error) {
	path := c.apiURL + "/market/stocks"

	return c.instruments(path)
}

func (c *RestClient) instruments(path string) ([]Instrument, error) {
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
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

func (c *RestClient) Operations(from time.Time, interval OperationInterval, figi string) ([]Operation, error) {
	q := url.Values{
		"from":     []string{from.Format(time.RFC3339)},
		"interval": []string{string(interval)},
	}
	if figi != "" {
		q.Set("figi", figi)
	}

	path := c.apiURL + "/operations?" + q.Encode()

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
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
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload.Operations, nil
}

func (c *RestClient) Portfolio() (Portfolio, error) {
	positions, err := c.PositionsPortfolio()
	if err != nil {
		return Portfolio{}, err
	}

	currencies, err := c.CurrenciesPortfolio()
	if err != nil {
		return Portfolio{}, err
	}

	return Portfolio{
		Currencies: currencies,
		Positions:  positions,
	}, nil
}

func (c *RestClient) PositionsPortfolio() ([]PositionBalance, error) {
	path := c.apiURL + "/portfolio"

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
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
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload.Positions, nil
}

func (c *RestClient) CurrenciesPortfolio() ([]CurrencyBalance, error) {
	path := c.apiURL + "/portfolio/currencies"

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
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
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload.Currencies, nil
}

func (c *RestClient) OrderCancel(id string) error {
	path := c.apiURL + "/orders/cancel?orderId=" + id

	return c.postJSONThrow(path, nil)
}

func (c *RestClient) LimitOrder(figi string, lots int, operation OperationType, price float64) (PlacedLimitOrder, error) {
	path := c.apiURL + "/orders/limit-order?figi=" + figi

	payload := struct {
		Lots      int           `json:"lots"`
		Operation OperationType `json:"operation"`
		Price     float64       `json:"price"`
	}{Lots: lots, Operation: operation, Price: price}

	bb, err := json.Marshal(payload)
	if err != nil {
		return PlacedLimitOrder{}, errors.Errorf("can't marshal request to %s body=%+v", path, payload)
	}

	req, err := c.newRequest(http.MethodPost, path, bytes.NewReader(bb))
	if err != nil {
		return PlacedLimitOrder{}, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return PlacedLimitOrder{}, err
	}

	type response struct {
		Payload PlacedLimitOrder `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return PlacedLimitOrder{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) Orders() ([]Order, error) {
	path := c.apiURL + "/orders"

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	type response struct {
		Payload []Order `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) Candles(from, to time.Time, interval CandleInterval, figi string) ([]Candle, error) {
	q := url.Values{
		"from":     []string{from.Format(time.RFC3339)},
		"to":       []string{to.Format(time.RFC3339)},
		"interval": []string{string(interval)},
		"figi":     []string{figi},
	}
	path := c.apiURL + "/market/candles?" + q.Encode()

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
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
		return nil, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload.Candles, nil
}

func (c *RestClient) Orderbook(depth int, figi string) (RestOrderBook, error) {
	if depth < 1 || depth > MaxOrderbookDepth {
		return RestOrderBook{}, ErrDepth
	}

	q := url.Values{
		"depth": []string{strconv.Itoa(depth)},
		"figi":  []string{figi},
	}
	path := c.apiURL + "/market/orderbook?" + q.Encode()

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return RestOrderBook{}, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return RestOrderBook{}, err
	}

	type response struct {
		Payload RestOrderBook `json:"payload"`
	}

	var resp response
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return RestOrderBook{}, errors.Wrapf(err, "can't unmarshal response to %s, respBody=%s", path, respBody)
	}

	return resp.Payload, nil
}

func (c *RestClient) postJSONThrow(url string, body interface{}) error {
	var bb []byte
	var err error

	if body != nil {
		bb, err = json.Marshal(body)
		if err != nil {
			return errors.Errorf("can't marshal request body to %s", url)
		}
	}

	req, err := c.newRequest(http.MethodPost, url, bytes.NewReader(bb))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *RestClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Errorf("can't create http request to %s", url)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	return req, nil
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
		return nil, errors.Errorf("bad response to %s code=%d, body=%s", req.URL.RawPath, resp.StatusCode, body)
	}

	return body, nil
}

type TradingError struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Payload    struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"payload"`
}

func (t TradingError) Error() string {
	return fmt.Sprintf(
		"TrackingID: %s, Status: %s, Message: %s, Code: %s",
		t.TrackingID, t.Status, t.Payload.Message, t.Payload.Code,
	)
}

func (t TradingError) NotEnoughBalance() bool {
	return t.Payload.Code == "NOT_ENOUGH_BALANCE"
}
