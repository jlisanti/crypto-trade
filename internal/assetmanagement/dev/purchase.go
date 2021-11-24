package assetmanagement

import (
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"strings"
)

type Trade struct {
	BuyCurrency  string
	SellCurrency string
	BuyValue     string
	SellValue    string
	BuyQuantity  string
	SellQuantity string
	ExchangeRate string
	Fees         string
	Exchange     string
}

func purchase(trade) {
	currency := BuyCurrecy + SellCurrency
	order := coinbasepro.Order{
		Price:     trade.BuyValue,
		Size:      trade.BuyQuantity,
		Side:      "buy",
		ProductID: "BTC-USD",
	}
}
