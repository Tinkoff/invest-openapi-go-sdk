package sdk

import "fmt"

// TradingError contains error info from tinkoff invest.
type TradingError struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Hint       string
	Payload    struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"payload"`
}

// Error for implements error.
func (t TradingError) Error() string {
	return fmt.Sprintf(
		"TrackingID: %s, Status: %s, Message: %s, Code: %s, Hint: %s",
		t.TrackingID, t.Status, t.Payload.Message, t.Payload.Code, t.Hint,
	)
}

// NotEnoughBalance for check error.
func (t TradingError) NotEnoughBalance() bool {
	return t.Payload.Code == "NOT_ENOUGH_BALANCE"
}

// InvalidTokenSpace for check error.
func (t TradingError) InvalidTokenSpace() bool {
	return t.Payload.Message == "Invalid token scopes"
}
