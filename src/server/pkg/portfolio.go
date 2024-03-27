package pkg

import (
	"errors"

	"github.com/google/uuid"
)

type Ticker string

const (
	TSLA Ticker = "tesla"
	META Ticker = "meta"
	NVDA Ticker = "nvidia"
	GOOG Ticker = "google"
)

type Portfolio struct {
	Cash   Currency
	Stocks map[Ticker]*Equity
	/* user to order map */
	History map[uuid.UUID]*OrderDto
	Id      uuid.UUID
}

type Equity struct {
	Ticker   Ticker
	Currency Currency
	Quantity int64
}

func NewPortfolio() *Portfolio {
	return &Portfolio{
		Cash:    Currency(100 * Eur),
		Stocks:  make(map[Ticker]*Equity),
		History: make(map[uuid.UUID]*OrderDto),
	}
}

func (p *Portfolio) isFulfillable(order *OrderDto) bool {
	if order.OrderType == Buy {
		return p.canBuy(Equity{
			Ticker:   order.Ticker,
			Currency: order.Price,
			Quantity: order.Quantity,
		})
	}

	if order.OrderType == Sell {
		return p.canSell(Equity{
			Ticker:   order.Ticker,
			Currency: order.Price,
			Quantity: order.Quantity,
		})
	}

	return false
}

func (p *Portfolio) canBuy(asset Equity) bool {
	return int64(p.Cash) > int64(asset.Quantity)*int64(asset.Currency)
}

func (p *Portfolio) canSell(asset Equity) bool {
	return p.Stocks[asset.Ticker].Quantity <= asset.Quantity
}

func (p *Portfolio) TryMakeTrade(o *OrderDto) (*OrderDto, error) {
	if p.isFulfillable(o) {
		p.History[o.Id] = o
		p.Cash -= o.Price * Currency(o.Quantity)
		return o, nil
	}

	return o, errors.New("not enough funds")
}
