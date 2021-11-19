package finance

import (
	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	"github.com/shopspring/decimal"
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
}

func computeROI(trade Trade, asset Asset) (roi float64, value string, age string) {
	FinalValueOfInvestment = decimal.NewFromString(trade.SellValue)
	InitialValueOfInvestment = decimal.NewFromString(asset.Quantity)
	*decimal.NewFromString(asset.BuyPrice)
	CostOfInvestment = decimal.NewFromString(trade.Fees)
	roi = ((FinalValueOfInvestment - InitialValueOfInvestment) / CostOfInvestment) * 100.0
	age = asset.BuyDate
	value = FinalValueOfInvestment - InitialValueOfInvestment
}
