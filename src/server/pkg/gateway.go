package pkg

import (
	"net"
	"server_exchange/shared/bindings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Gateway struct {
	Conn net.Conn
}

type ClientHandler interface {
	Handle(acc *Portfolio, op *OrderProcessor)

	Buy(req *bindings.ServerRequest) (*OrderDto, error)
	Sell(req *bindings.ServerRequest) (*OrderDto, error)
	List(req *bindings.ServerRequest) (*map[Ticker]*Equity, error)
	Ok(res *bindings.Response) *bindings.ServerResponse
	Ko(err error) *bindings.ServerResponse
}

func (g *Gateway) Handle(acc *Portfolio, op *OrderProcessor) {
	defer g.Conn.Close()

	for {

		buf := make([]byte, 1024)
		serverRequest := bindings.ServerRequest{}
		err := proto.Unmarshal(buf, &serverRequest)

		if err != nil {
			order, err := g.Buy(&serverRequest)
			var res *bindings.ServerResponse = nil
			if err != nil {
				_, err := acc.TryMakeTrade(order)
				if err == nil {
					op.Ingress <- &CustomerOrder{
						CustomerId: acc.Id,
						OrderDto:   order,
					}
					res = g.Ok(&bindings.Response{
						Response: &bindings.Response_Order{
							Order: &bindings.OrderResponse{
								Id: order.Id.String(),
							},
						},
					})

				}
			} else {
				res = g.Ko(err)
			}

			data, err := proto.Marshal(res)
			if err != nil {
				g.Conn.Write(data)
			}
		}
	}
}
func (c *Gateway) Buy(req *bindings.ServerRequest) (*OrderDto, error) {
	order := req.GetOrder()
	if order != nil {
		return &OrderDto{
			Id:        uuid.New(),
			Quantity:  int64(order.Quantity),
			Ticker:    Ticker(order.Ticker),
			Price:     ToCurrency(order.Price * float32(Cent)),
			OrderType: OrderType(order.OrderType),
		}, nil
	} else {
		return nil, &net.ParseError{}
	}
}

func (c *Gateway) Sell(req *bindings.ServerRequest) (*OrderDto, error) {
	return nil, &net.ParseError{}
}

func (c *Gateway) List(req *bindings.ServerRequest) (*map[Ticker]*Equity, error) {
	return nil, &net.ParseError{}
}

func (g *Gateway) Ko(err error) *bindings.ServerResponse {
	return &bindings.ServerResponse{
		WasSuccessful: false,
		Result: &bindings.ServerResponse_Error{
			Error: &bindings.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		},
	}
}

func (g *Gateway) Ok(res *bindings.Response) *bindings.ServerResponse {
	return &bindings.ServerResponse{
		WasSuccessful: true,
		Result: &bindings.ServerResponse_Response{
			Response: res,
		},
	}
}
