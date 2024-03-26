package server

import (
	"sync"

	"github.com/google/btree"
	"github.com/google/uuid"
)

type LimitBook struct {
	Asks           map[Ticker]*btree.BTree
	Bids           map[Ticker]*btree.BTree
	WriteLocks     map[Ticker]*sync.Mutex
	OrderProcessor *OrderProcessor
}

type LimitTrade struct {
	Id    uuid.UUID
	Trade *Trade
}

type LimitTradingProcessing interface {
	InsertBuy(*LimitTrade) uuid.UUID
	InsertSell(*LimitTrade) uuid.UUID
	FillBid(*LimitTrade)
	FillAsk(*LimitTrade)
	PartialFill(*LimitTrade, int64)
}

func NewLimitBook(op *OrderProcessor, tickers []Ticker) *LimitBook {
	lBook := LimitBook{
		Asks:       make(map[Ticker]*btree.BTree),
		Bids:       make(map[Ticker]*btree.BTree),
		WriteLocks: make(map[Ticker]*sync.Mutex),
	}

	for i := 0; i < len(tickers); i++ {
		lBook.Asks[tickers[i]] = btree.New(2)
		lBook.Bids[tickers[i]] = btree.New(2)
		lBook.WriteLocks[tickers[i]] = &sync.Mutex{}
	}
	return &lBook
}

func (lTrade LimitTrade) Less(than btree.Item) bool {
	upcast := than.(LimitTrade)

	return lTrade.Id != upcast.Id &&
		lTrade.Trade.Price < upcast.Trade.Price
}

func (lBook LimitBook) InsertBuy(lTrade *LimitTrade) {
	ticker := lTrade.Trade.Ticker

	lBook.WriteLocks[ticker].Lock()
	defer lBook.WriteLocks[ticker].Unlock()

	lBook.Bids[ticker].ReplaceOrInsert(LimitTrade{Id: lTrade.Id, Trade: lTrade.Trade})
}

func (lBook LimitBook) InsertSell(lTrade *LimitTrade) {
	ticker := lTrade.Trade.Ticker

	lBook.WriteLocks[ticker].Lock()
	defer lBook.WriteLocks[ticker].Unlock()

	lBook.Asks[ticker].ReplaceOrInsert(LimitTrade{Id: lTrade.Id, Trade: lTrade.Trade})
}

func (l LimitBook) Init(t []Ticker) {

}

func (book *LimitBook) FillBid(t LimitTrade) {
	ticker := t.Trade.Ticker
	book.WriteLocks[ticker].Lock()
	defer book.WriteLocks[ticker].Unlock()
	book.Bids[ticker].Delete(t)
	book.OrderProcessor.Egress <- &t
}

func (book *LimitBook) FillAsk(t LimitTrade) {
	ticker := t.Trade.Ticker
	book.WriteLocks[ticker].Lock()
	defer book.WriteLocks[ticker].Unlock()
	book.Asks[ticker].Delete(t)
	book.OrderProcessor.Egress <- &t
}

func (book *LimitBook) PartialFill(t LimitTrade, taken int64) {
	ticker := t.Trade.Ticker
	book.WriteLocks[ticker].Lock()
	defer book.WriteLocks[ticker].Unlock()
	t.Trade.Quantity -= taken
}
