package client

import (
	"fmt"
	"net"

	"github.com/alecthomas/kong"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:9032")
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := kong.Parse(&CLI)
	server := &Gateway{Conn: conn}

	for {
		switch ctx.Command() {
		case "buy <ticker> <price> <quantity>":
			fmt.Println(server.Buy(&CLI))
		case "sell <ticker> <price> <quantity>":
			fmt.Println(server.Sell(&CLI))
		case "list":
			fmt.Println(server.List(&CLI))
		default:
			panic(ctx.Command())
		}
	}
}
