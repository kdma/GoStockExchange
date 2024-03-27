package pkg

import "github.com/google/uuid"

type CustomerOrder struct {
	CustomerId uuid.UUID
	OrderDto   *OrderDto
}
type OrderDto struct {
	Id        uuid.UUID
	Ticker    Ticker
	Price     Currency
	Quantity  int64
	OrderType OrderType
}
type OrderType int

const (
	Buy  OrderType = 0
	Sell OrderType = 1
)
