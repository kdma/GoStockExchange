package pkg

import "github.com/google/uuid"

type CustomerOperations struct {
	Id         uuid.UUID
	Operations map[uuid.UUID]*StatefulOrder
}

type OrderStatus int

const (
	Created   OrderStatus = 0
	Cancelled OrderStatus = 1
	Fulfilled OrderStatus = 2
)

type StatefulOrder struct {
	Status OrderStatus
	Trade  *LimitTrade
}

func (co *CustomerOperations) Add(s *StatefulOrder) {
	co.Operations[s.Trade.Id] = s
}

func (co *CustomerOperations) Mutate(id uuid.UUID, status OrderStatus) {
	co.Operations[id].Status = status
}
