package main

import (
	"fmt"
	"net"
	"server_exchange/server/pkg"

	"server_exchange/shared/bindings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func main() {

	l, err := net.Listen("tcp4", "localhost:9032")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	tickers := []pkg.Ticker{pkg.TSLA, pkg.NVDA, pkg.META, pkg.GOOG}

	treasury := pkg.NewTreasury()
	orderProcessor := pkg.NewOrderProcessor()
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

func clientHandler(c net.Conn, acc *pkg.Portfolio, op *pkg.OrderProcessor) {
	defer c.Close()

	for {

		buf := make([]byte, 1024)
		serverRequest := bindings.ServerRequest{}
		err := proto.Unmarshal(buf, &serverRequest)

		if err != nil {
			order := serverRequest.GetOrder()
			if order != nil {
				orderDto := &pkg.OrderDto{
					Id:        uuid.New(),
					Quantity:  int64(order.Quantity),
					Ticker:    pkg.Ticker(order.Ticker),
					Price:     pkg.Currency(order.Price),
					OrderType: pkg.OrderType(order.OrderType),
				}

				_, err := acc.TryMakeTrade(orderDto)
				if err == nil {
					op.Ingress <- &pkg.CustomerOrder{
						CustomerId: acc.Id,
						OrderDto:   orderDto,
					}
				}
				//send a response
			}

		}
	}
}
