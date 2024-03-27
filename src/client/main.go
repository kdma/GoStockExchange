package main

import (
	"encoding/json"
	"fmt"
	"net"
	"server_exchange/client/pkg"
	"server_exchange/shared/bindings"
	"strings"

	"github.com/alecthomas/kong"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:9032")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	server := &pkg.Gateway{Conn: conn}

	for {
		err, ctx := parse(pkg.CLI)
		if err != nil {
			fmt.Println(err.Error())
		}

		switch ctx.Command() {
		case "mkt":
			handleResponse(server.Market(&pkg.CLI))
		case "buy <ticker> <price> <quantity>":
			handleResponse(server.Buy(&pkg.CLI))
		case "sell <ticker> <price> <quantity>":
			handleResponse(server.Sell(&pkg.CLI))
		case "list":
			handleResponse(server.List(&pkg.CLI))
		default:
		}
	}
}

func handleResponse(e error, res *bindings.ServerResponse) {
	if e != nil {
		fmt.Println(e.Error())
	} else {
		if res.WasSuccessful {
			fmt.Print(json.Marshal(*res.GetData()))
		} else {
			fmt.Print(json.Marshal(*res.GetError()))
		}
	}
}

func parse(cli interface{}, options ...kong.Option) (error, *kong.Context) {
	parser, err := kong.New(cli, options...)
	if err != nil {
		panic(err)
	}

	var line string
	_, err = fmt.Scanln(&line)
	if err != nil {
		ctx, err := parser.Parse(strings.Split(line, " "))
		parser.FatalIfErrorf(err)
		return nil, ctx
	}

	return err, nil

}
