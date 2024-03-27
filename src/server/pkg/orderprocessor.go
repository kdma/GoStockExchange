package pkg

import "github.com/google/uuid"

type OrderProcessor struct {
	Ingress            chan *CustomerOrder
	Egress             chan *LimitTrade
	Operations         map[uuid.UUID]*CustomerOperations
	CustomerOperations map[uuid.UUID]uuid.UUID
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		Ingress:            make(chan *CustomerOrder),
		Egress:             make(chan *LimitTrade),
		Operations:         make(map[uuid.UUID]*CustomerOperations),
		CustomerOperations: make(map[uuid.UUID]uuid.UUID),
	}
}

func (op *OrderProcessor) Process(tickers []Ticker) {

	book := NewLimitBook(op, tickers)

	go op.ListenNewOrders(book)
	go op.ListenFulfilled()

	marketMaker := &MarketMaker{LimitBook: book}
	marketMaker.Match(tickers)

}
func (op *OrderProcessor) ListenNewOrders(book *LimitBook) {
	for {

		reqOrder := <-op.Ingress
		trade := &LimitTrade{
			Id:    uuid.New(),
			Trade: &Trade{Ticker: reqOrder.OrderDto.Ticker, Price: reqOrder.OrderDto.Price, Quantity: reqOrder.OrderDto.Quantity, OrderType: reqOrder.OrderDto.OrderType},
		}

		if op.Operations[reqOrder.CustomerId] == nil {
			op.Operations[reqOrder.CustomerId] = &CustomerOperations{}
		}

		op.CustomerOperations[trade.Id] = reqOrder.CustomerId
		op.Operations[reqOrder.CustomerId].Add(&StatefulOrder{Status: Created, Trade: trade})

		if trade.Trade.OrderType == Buy {
			book.InsertBuy(trade)
		}
		if trade.Trade.OrderType == Sell {
			book.InsertSell(trade)
		}
	}
}

func (op *OrderProcessor) ListenFulfilled() {
	for {
		fulfilled := <-op.Egress

		custId := op.CustomerOperations[fulfilled.Id]
		op.Operations[custId].Mutate(fulfilled.Id, Fulfilled)
	}
}
