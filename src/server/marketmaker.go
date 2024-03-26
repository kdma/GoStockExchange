package server

type MarketMaker struct {
	LimitBook *LimitBook
}

type MarketMaking interface {
	Match()
	Trade()
}

func (m *MarketMaker) Match(tickers []Ticker) {

	for _, t := range tickers {
		go m.trade(t)
	}
}

func (m *MarketMaker) trade(t Ticker) {

	for {

		bids := m.LimitBook.Bids[t]
		asks := m.LimitBook.Asks[t]

		maxBid := bids.Max().(LimitTrade)
		minAsk := asks.Min().(LimitTrade)

		if maxBid.Trade.Price >= minAsk.Trade.Price {

			if maxBid.Trade.Quantity == minAsk.Trade.Quantity {

				m.LimitBook.FillBid(maxBid)
				m.LimitBook.FillAsk(minAsk)

			} else if maxBid.Trade.Quantity < minAsk.Trade.Quantity {

				m.LimitBook.FillBid(maxBid)
				m.LimitBook.PartialFill(minAsk, minAsk.Trade.Quantity)

			} else if maxBid.Trade.Quantity > minAsk.Trade.Quantity {

				m.LimitBook.PartialFill(maxBid, minAsk.Trade.Quantity)
				m.LimitBook.FillAsk(minAsk)
			}
		}

	}
}
