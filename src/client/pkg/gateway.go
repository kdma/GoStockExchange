package pkg

import (
	"log"
	"net"
	"server_exchange/shared/bindings"

	"github.com/golang/protobuf/proto"
)

type Gateway struct {
	Conn net.Conn
}
type Server interface {
	Buy(c *Cli) (error, *bindings.ServerResponse)
	Sell(c *Cli) (error, *bindings.ServerResponse)
	List(c *Cli) (error, *bindings.ServerResponse)
	Market(c *Cli) (error, *bindings.ServerResponse)
}

func readResponse(c net.Conn) (error, *bindings.ServerResponse) {
	recvBuf := make([]byte, 1024)

	_, err := c.Read(recvBuf[:]) // recv data
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return netErr, nil
			// time out
		} else {
			log.Println("read error:", err)
			return err, nil
			// some error else, do something else, for example create new conn
		}
	}

	res := &bindings.ServerResponse{}
	err = proto.Unmarshal(recvBuf, res)
	if err != nil {
		return err, nil
	}
	return nil, res
}

func (g *Gateway) Buy(c *Cli) (error, *bindings.ServerResponse) {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_Order{
			Order: &bindings.OrderRequest{
				Price:     CLI.Buy.Price,
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

	return readResponse(g.Conn)
}

func (g *Gateway) Sell(c *Cli) (error, *bindings.ServerResponse) {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_Order{
			Order: &bindings.OrderRequest{
				Price:     CLI.Buy.Price,
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

	return readResponse(g.Conn)
}

func (g *Gateway) List(c *Cli) (error, *bindings.ServerResponse) {
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

	return readResponse(g.Conn)
}

func (g *Gateway) Market(c *Cli) (error, *bindings.ServerResponse) {
	req := &bindings.ServerRequest{
		Request: &bindings.ServerRequest_Tickers{
			Tickers: &bindings.MarketRequest{},
		},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		g.Conn.Write(data)
	}

	return readResponse(g.Conn)
}
