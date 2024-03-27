package pkg

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

		var maxBid *LimitTrade
		var minAsk *LimitTrade

		if bids.Max() != nil {
			bid := bids.Max().(LimitTrade)
			maxBid = &bid
		}

		if bids.Min() != nil {
			bid := asks.Min().(LimitTrade)
			minAsk = &bid
		}

		if maxBid != nil && minAsk != nil && maxBid.Trade.Price >= minAsk.Trade.Price {

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
