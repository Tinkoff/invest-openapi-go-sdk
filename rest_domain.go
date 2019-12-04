package sdk

import "time"

type Currency string

const (
	RUB Currency = "RUB"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

type OperationType string

const (
	BUY  OperationType = "Buy"
	SELL OperationType = "Sell"
)

type OrderStatus string

const (
	OrderStatusNew            OrderStatus = "New"
	OrderStatusPartiallyFill  OrderStatus = "PartiallyFill"
	OrderStatusFill           OrderStatus = "Fill"
	OrderStatusCancelled      OrderStatus = "Cancelled"
	OrderStatusReplaced       OrderStatus = "Replaced"
	OrderStatusPendingCancel  OrderStatus = "PendingCancel"
	OrderStatusRejected       OrderStatus = "Rejected"
	OrderStatusPendingReplace OrderStatus = "PendingReplace"
	OrderStatusPendingNew     OrderStatus = "PendingNew"
)

type OperationStatus string

const (
	OperationStatusDone     OrderStatus = "Done"
	OperationStatusDecline  OrderStatus = "Decline"
	OperationStatusProgress OrderStatus = "Progress"
)

type PlacedLimitOrder struct {
	ID            string      `json:"orderId"`
	Operation     string      `json:"operation"`
	Status        OrderStatus `json:"status"`
	RejectReason  string      `json:"rejectReason"`
	RequestedLots int         `json:"requestedLots"`
	ExecutedLots  int         `json:"executedLots"`
	Commission    MoneyAmount `json:"commission"`
}

type Order struct {
	ID            string  `json:"orderId"`
	FIGI          string  `json:"figi"`
	Operation     string  `json:"operation"`
	Status        string  `json:"status"`
	RequestedLots int     `json:"requestedLots"`
	ExecutedLots  int     `json:"requestedLots"`
	Type          string  `json:"type"`
	Price         float64 `json:"price"`
}

type Portfolio struct {
	Positions  []PositionBalance
	Currencies []CurrencyBalance
}

type CurrencyBalance struct {
	Currency Currency `json:"currency"`
	Balance  float64  `json:"balance"`
	Blocked  float64  `json:"blocked"`
}

type PositionBalance struct {
	FIGI           string      `json:"figi"`
	Ticker         string      `json:"ticker"`
	ISIN           string      `json:"isin"`
	InstrumentType string      `json:"instrumentType"`
	Balance        float64     `json:"balance"`
	Blocked        float64     `json:"blocked"`
	Lots           int         `json:"lots"`
	ExpectedYield  MoneyAmount `json:"expectedYield"`
}

type MoneyAmount struct {
	Currency Currency `json:"currency"`
	Value    float64  `json:"value"`
}

type Instrument struct {
	FIGI              string   `json:"figi"`
	Ticker            string   `json:"ticker"`
	ISIN              string   `json:"isin"`
	Name              string   `json:"name"`
	MinPriceIncrement float64  `json:"minPriceIncrement"`
	Lot               int      `json:"lot"`
	Currency          Currency `json:"currency"`
}

type Operation struct {
	ID             string          `json:"id"`
	Status         OperationStatus `json:"status"`
	Trades         []Trade         `json:"trades"`
	Commission     MoneyAmount     `json:"commission"`
	Currency       Currency        `json:"currency"`
	Payment        float64         `json:"payment"`
	Price          float64         `json:"price"`
	Quantity       int             `json:"quantity"`
	FIGI           string          `json:"figi"`
	InstrumentType string          `json:"instrumentType"`
	IsMarginCall   bool            `json:"isMarginCall"`
	DateTime       time.Time       `json:"date"`
	OperationType  OperationType   `json:"operationType"`
}

type Trade struct {
	ID       string    `json:"tradeId"`
	DateTime time.Time `json:"date"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
}

type RestPriceQuantity struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type RestOrderBook struct {
	FIGI              string              `json:"figi"`
	Depth             int                 `json:"depth"`
	Bids              []RestPriceQuantity `json:"bids"`
	Asks              []RestPriceQuantity `json:"asks"`
	TradeStatus       TradingStatus       `json:"tradeStatus"`
	MinPriceIncrement float64             `json:"minPriceIncrement"`
	LastPrice         float64             `json:"lastPrice,omitempty"`
	ClosePrice        float64             `json:"closePrice,omitempty"`
	LimitUp           float64             `json:"limitUp,omitempty"`
	LimitDown         float64             `json:"limitDown,omitempty"`
}
