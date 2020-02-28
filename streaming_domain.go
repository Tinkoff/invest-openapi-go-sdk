package sdk

import "time"

const MaxOrderbookDepth = 20

type CandleInterval string

const (
	CandleInterval1Min   CandleInterval = "1min"
	CandleInterval2Min   CandleInterval = "2min"
	CandleInterval3Min   CandleInterval = "3min"
	CandleInterval5Min   CandleInterval = "5min"
	CandleInterval10Min  CandleInterval = "10min"
	CandleInterval15Min  CandleInterval = "15min"
	CandleInterval30Min  CandleInterval = "30min"
	CandleInterval1Hour  CandleInterval = "hour"
	CandleInterval2Hour  CandleInterval = "2hour"
	CandleInterval4Hour  CandleInterval = "4hour"
	CandleInterval1Day   CandleInterval = "day"
	CandleInterval1Week  CandleInterval = "week"
	CandleInterval1Month CandleInterval = "month"
)

type TradingStatus string

const (
	BreakInTrading               TradingStatus = "break_in_trading"
	NormalTrading                TradingStatus = "normal_trading"
	NotAvailableForTrading       TradingStatus = "not_available_for_trading"
	ClosingAuction               TradingStatus = "closing_auction"
	ClosingPeriod                TradingStatus = "closing_period"
	DarkPoolAuction              TradingStatus = "dark_pool_auction"
	DiscreteAuction              TradingStatus = "discrete_auction"
	OpeningPeriod                TradingStatus = "opening_period"
	OpeningAuctionPeriod         TradingStatus = "opening_auction_period"
	TradingAtClosingAuctionPrice TradingStatus = "trading_at_closing_auction_price"
)

type Event struct {
	Name string `json:"event"`
}

type FullEvent struct {
	Name string    `json:"event"`
	Time time.Time `json:"time"`
}

type CandleEvent struct {
	FullEvent
	Candle Candle `json:"payload"`
}

type Candle struct {
	FIGI       string         `json:"figi"`
	Interval   CandleInterval `json:"interval"`
	OpenPrice  float64        `json:"o"`
	ClosePrice float64        `json:"c"`
	HighPrice  float64        `json:"h"`
	LowPrice   float64        `json:"l"`
	Volume     float64        `json:"v"`
	TS         time.Time      `json:"time"`
}

type OrderBookEvent struct {
	FullEvent
	OrderBook OrderBook `json:"payload"`
}

type OrderBook struct {
	FIGI  string          `json:"figi"`
	Depth int             `json:"depth"`
	Bids  []PriceQuantity `json:"bids"`
	Asks  []PriceQuantity `json:"asks"`
}

type PriceQuantity [2]float64 // 0 - price, 1 - quantity

type InstrumentInfoEvent struct {
	FullEvent
	Info InstrumentInfo `json:"payload"`
}

type InstrumentInfo struct {
	FIGI              string        `json:"figi"`
	TradeStatus       TradingStatus `json:"trade_status"`
	MinPriceIncrement float64       `json:"min_price_increment"`
	Lot               float64       `json:"lot"`
	AccruedInterest   float64       `json:"accrued_interest,omitempty"`
	LimitUp           float64       `json:"limit_up,omitempty"`
	LimitDown         float64       `json:"limit_down,omitempty"`
}

type ErrorEvent struct {
	FullEvent
	Error Error `json:"payload"`
}

type Error struct {
	RequestID string `json:"request_id,omitempty"`
	Error     string `json:"error"`
}
