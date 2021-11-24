package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	//"os"
	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	//"github.com/shopspring/decimal"
	//"github.com/jlisanti/crypto-trade/internal/finance"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

func main() {

	// configure coinbasepro

	client := coinbasepro.NewClient()

	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
		Key:        "89b1f52167924567c1a41b42a236d8a1",
		Passphrase: "puj31du7a4j",
		Secret:     "RUWPjira048friEd52Z34ptpYdeFnop1PrucxrvGRlZhUtNuM71Iub+HwTu7X2Gg8OjVkuFIW1iPm5C8qzamgw==",
	})

	accounts, err := client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	var assets = []assetmanagement.Asset{}

	// Pull asset record from coinbasepro

	var ledgers []coinbasepro.LedgerEntry

	for _, a := range accounts {
		cursor := client.ListAccountLedger(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledgers); err != nil {
				println(err.Error())
			}
			for _, e := range ledgers {
				if e.Type != "transfer" {
					currencies := strings.Split(e.Details.ProductID, "-")

					// Determine if this was a buy or sell
					transferAmount, _ := strconv.ParseFloat(e.Amount, 64)
					if transferAmount > float64(0.0) {

						// Store asset
						asset := assetmanagement.Asset{
							Index:    e.ID,
							Currency: currencies[0],
							Quantity: e.Amount,
							BuyDate:  time.Time(e.CreatedAt.Time()),
							BuyPrice: "",
							Cost:     "",
						}
						assets = append(assets, asset)

					}
				}
			}
		}
	}

	fmt.Println("Assets \n")
	for i := range assets {
		fmt.Println(assets[i].Currency)
	}

}
