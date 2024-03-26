package server

type Trade struct {
	Ticker    Ticker
	Price     Currency
	Quantity  int64
	OrderType OrderType
}

func fromOrderDto(o *OrderDto) *Trade {
	return &Trade{
		Ticker:    o.Ticker,
		Price:     o.Price,
		Quantity:  o.Quantity,
		OrderType: o.OrderType,
	}
}
