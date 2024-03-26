package server

import (
	"fmt"
	"net"
	"server_exchange/shared/bindings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func main() {

	l, err := net.Listen("tcp4", "9032")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	tickers := []Ticker{TSLA, NVDA, META, GOOG}

	treasury := NewTreasury()
	orderProcessor := NewOrderProcessor()
	orderProcessor.Process(tickers)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		uId := uuid.New()
		treasury.AddCustomer(uId)
		go clientHandler(c, treasury.GetCustomer(uId), orderProcessor)
	}
}

func clientHandler(c net.Conn, acc *Portfolio, op *OrderProcessor) {
	defer c.Close()

	for {

		buf := make([]byte, 1024)
		serverRequest := bindings.ServerRequest{}
		err := proto.Unmarshal(buf, &serverRequest)

		if err != nil {
			order := serverRequest.GetOrder()
			if order != nil {
				orderDto := &OrderDto{
					Id:        uuid.New(),
					Quantity:  int64(order.Quantity),
					Ticker:    Ticker(order.Ticker),
					Price:     Currency(order.Price),
					OrderType: OrderType(order.OrderType),
				}

				_, err := acc.TryMakeTrade(orderDto)
				if err == nil {
					op.Ingress <- &CustomerOrder{
						CustomerId: acc.Id,
						OrderDto:   orderDto,
					}
				}
			}

		}
	}
}

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
