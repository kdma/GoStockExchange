package client

import (
	"net"
	"server_exchange/shared/bindings"

	"github.com/golang/protobuf/proto"
)

type Gateway struct {
	Conn net.Conn
}
type Server interface {
	Buy(c *Cli) *bindings.ServerResponse
	Sell(c *Cli) *bindings.ServerResponse
	List(c *Cli) *bindings.ServerResponse
}

func (g *Gateway) Buy(c *Cli) *bindings.ServerResponse {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_Order{
			Order: &bindings.OrderRequest{
				Price:     int32(CLI.Buy.Price),
				Ticker:    CLI.Buy.Ticker,
				OrderType: 0,
				Quantity:  CLI.Buy.Quantity,
			},
		},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		g.Conn.Write(data)
	}

	return &bindings.ServerResponse{}
}

func (g *Gateway) Sell(c *Cli) *bindings.ServerResponse {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_Order{
			Order: &bindings.OrderRequest{
				Price:     int32(CLI.Buy.Price),
				Ticker:    CLI.Buy.Ticker,
				OrderType: 1,
				Quantity:  CLI.Buy.Quantity,
			},
		},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		g.Conn.Write(data)
	}

	return &bindings.ServerResponse{}
}

func (g *Gateway) List(c *Cli) *bindings.ServerResponse {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_List{
			List: &bindings.ListRequest{
				User: "",
			},
		},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		g.Conn.Write(data)
	}

	return &bindings.ServerResponse{}
}
