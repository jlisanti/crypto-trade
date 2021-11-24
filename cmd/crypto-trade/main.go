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
	//    need to test that asset information is correctly populated
	//    from the coinbase record

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
							ID:       e.ID,
							Currency: currencies[0],
							Quantity: e.Amount,
							BuyDate:  time.Time(e.CreatedAt.Time()),
							BuyPrice: "",
							Cost:     "",
						}
						assets = append(assets, asset)

						// Store purchase price
					} else {
						// Do I need a loop?
						for _, asset := range assets {
							if asset.ID == e.ID {
								asset.BuyPrice = e.Amount
							}
						}
					}

					// Store fee
				} else if e.Type == "fee" {
					for _, asset := range assets {
						if asset.ID == e.ID {
							asset.Cost = e.Amount
						}
					}
				}
			}
		}
	}

	fmt.Println("Assets \n")
	for _, asset := range assets {
		fmt.Println(asset.Currency)
		fmt.Println(asset.Cost)
		fmt.Println(asset.BuyPrice)
	}

}
