package client

type BuyCommand struct {
	Ticker   string `help:"Ticker"`
	Price    int64  `help:"Price in eur"`
	Quantity int64  `help:"Quantity 1 * "`
}

type SellCommand struct {
	Ticker   string `help:"Ticker"`
	Price    int64  `help:"Price in eur"`
	Quantity int64  `help:"Quantity 1 * "`
}

type ListCommand struct {
}
type Cli struct {
	Buy BuyCommand `cmd:"" help:"Buy stock."`

	Sell SellCommand `cmd:"" help:"Sell stock."`

	List ListCommand `cmd:"" help:"List portfolio."`
}

var CLI Cli
