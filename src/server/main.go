package main

import (
	"fmt"
	"net"
	"server_exchange/server/pkg"

	"github.com/google/uuid"
)

func main() {

	l, err := net.Listen("tcp4", "localhost:9032")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	treasury := pkg.NewTreasury()
	orderProcessor := pkg.NewOrderProcessor()
	orderProcessor.Process([]pkg.Ticker{pkg.TSLA, pkg.NVDA, pkg.META, pkg.GOOG})

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		uId := uuid.New()
		treasury.AddCustomer(uId)
		gateway := &pkg.Gateway{
			Conn: c,
		}
		go gateway.Handle(treasury.GetCustomer(uId), orderProcessor)
	}
}
