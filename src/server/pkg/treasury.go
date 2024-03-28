package pkg

import (
	"github.com/google/uuid"
)

type Currency float32

func ToCurrency(cents float32) Currency {
	return Currency(cents) / Cent
}

const (
	Cent Currency = 100
	Eur  Currency = 10000
)

type Treasury struct {
	Accounts map[uuid.UUID]*Portfolio
}

type Operations interface {
	AddCustomer(uuid.UUID)
	GetCustomer(uuid.UUID) *Portfolio
}

func NewTreasury() *Treasury {
	return &Treasury{
		Accounts: make(map[uuid.UUID]*Portfolio),
	}
}

func (t *Treasury) AddCustomer(id uuid.UUID) {
	t.Accounts[id] = NewPortfolio()
}

func (t *Treasury) GetCustomer(id uuid.UUID) *Portfolio {
	return t.Accounts[id]
}
