package pkg

type BuyCommand struct {
	Ticker   string  `help:"Ticker"`
	Price    float32 `help:"Price in eur"`
	Quantity int32   `help:"Quantity 1 * "`
}

type SellCommand struct {
	Ticker   string  `help:"Ticker"`
	Price    float32 `help:"Price in eur"`
	Quantity int32   `help:"Quantity 1 * "`
}

type ListCommand struct {
}

type MarketCommand struct {
}
type Cli struct {
	Buy BuyCommand `cmd:"" help:"Buy stock."`

	Sell SellCommand `cmd:"" help:"Sell stock."`

	List ListCommand `cmd:"" help:"List portfolio."`

	Market MarketCommand `cmd:"" help:"List tickers."`
}

var CLI Cli
