package sdk

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const MarketApiURL = "wss://api-invest.tinkoff.ru/openapi/md/v1/md-openapi/ws"

type Logger interface {
	Printf(format string, args ...interface{})
}

type MarketClient struct {
	logger Logger
	conn   *websocket.Conn
	token  string
	apiURL string
}

func NewMarketClient(logger Logger, token string) (*MarketClient, error) {
	return NewMarketClientCustom(logger, token, MarketApiURL)
}

func NewMarketClientCustom(logger Logger, token, apiURL string) (*MarketClient, error) {
	client := &MarketClient{
		logger: logger,
		token:  token,
		apiURL: apiURL,
	}

	conn, err := client.connect()
	if err != nil {
		return nil, err
	}
	client.conn = conn

	return client, nil
}

func (c *MarketClient) Close() error {
	return c.conn.Close()
}

func (c *MarketClient) RunReadLoop(fn func(event interface{}) error) error {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "can't read message")
		}

		var event Event
		if err := json.Unmarshal(msg, &event); err != nil {
			c.logger.Printf("Can't unmarshal event %s", msg)
			continue
		}

		switch event.Name {
		case "candle":
			var event CandleEvent
			if err := json.Unmarshal(msg, &event); err != nil {
				c.logger.Printf("Can't unmarshal event candle %s", msg)
				continue
			}
			if err := fn(event); err != nil {
				return err
			}
		case "orderbook":
			var event OrderBookEvent
			if err := json.Unmarshal(msg, &event); err != nil {
				c.logger.Printf("Can't unmarshal event orderbook %s", msg)
				continue
			}
			if err := fn(event); err != nil {
				return err
			}
		case "instrument_info":
			var event InstrumentInfoEvent
			if err := json.Unmarshal(msg, &event); err != nil {
				c.logger.Printf("Can't unmarshal event instrument_info %s", msg)
				continue
			}
			if err := fn(event); err != nil {
				return err
			}
		case "error":
			var event ErrorEvent
			if err := json.Unmarshal(msg, &event); err != nil {
				c.logger.Printf("Can't unmarshal event error %s", msg)
				continue
			}
			if err := fn(event); err != nil {
				return err
			}
		default:
			c.logger.Printf("Get unknown event %s", msg)
		}
	}
}

func (c *MarketClient) SubscribeCandle(figi string, interval CandleInterval, requestID string) error {
	sub := `{ "event": "candle:subscribe", "request_id": "` + requestID + `", "figi": "` + figi + `", "interval": "` + string(interval) + `"}`

	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't subscribe to event")
	}

	return nil
}

func (c *MarketClient) UnsubscribeCandle(figi string, interval CandleInterval, requestID string) error {
	sub := `{ "event": "candle:unsubscribe", "request_id": "` + requestID + `", "figi": "` + figi + `", "interval": "` + string(interval) + `"}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't unsubscribe from event")
	}

	return nil
}

func (c *MarketClient) SubscribeOrderbook(figi string, depth int, requestID string) error {
	const maxDepth = 20
	if depth < 1 || depth > maxDepth {
		return errors.New("invalid depth. Should be in interval 0 < x <= 20")
	}

	sub := `{ "event": "orderbook:subscribe", "request_id": "` + requestID + `", "figi": "` + figi + `", "depth": ` + strconv.Itoa(depth) + `}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't subscribe to event")
	}

	return nil
}

func (c *MarketClient) UnsubscribeOrderbook(figi string, depth int, requestID string) error {
	const maxDepth = 20
	if depth < 1 || depth > maxDepth {
		return errors.New("invalid depth. Should be in interval 0 < x <= 20")
	}

	sub := `{ "event": "orderbook:unsubscribe", "request_id": "` + requestID + `", "figi": "` + figi + `", "depth": ` + strconv.Itoa(depth) + `}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't unsubscribe from event")
	}

	return nil
}

func (c *MarketClient) SubscribeInstrumentInfo(figi, requestID string) error {
	sub := `{"event": "instrument_info:subscribe", "request_id": "` + requestID + `", "figi": "` + figi + `"}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't subscribe to event")
	}

	return nil
}

func (c *MarketClient) UnsubscribeInstrumentInfo(figi, requestID string) error {
	sub := `{"event": "instrument_info:unsubscribe", "request_id": "` + requestID + `", "figi": "` + figi + `"}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sub)); err != nil {
		return errors.Wrap(err, "can't unsubscribe from event")
	}

	return nil
}

func (c *MarketClient) connect() (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}

	conn, resp, err := dialer.Dial(c.apiURL, http.Header{"Authorization": {"Bearer " + c.token}})
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to %s", c.apiURL)
	}
	defer resp.Body.Close()

	conn.SetPingHandler(func(message string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	return conn, nil
}
