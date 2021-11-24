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

	ws "github.com/gorilla/websocket"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

func main() {

	// configure coinbasepro

	client := coinbasepro.NewClient()
	/*

		client.UpdateConfig(&coinbasepro.ClientConfig{
			BaseURL:    "https://api.pro.coinbase.com/",
			Key:        "7ec721c05fdaa802a432775c565eeff9",
			Passphrase: "jr1o3qv98cf",
			Secret:     "o7Pr8LggG1ywJxsaLwQpmpEWLA87NjNr2cuX7caDw4GHMB86G0C+3hNu6eOSuIzkMKoBzYeKSQumClixXkPB3Q==",
		})
	*/

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

	println("account connection established")
	var assets = []assetmanagement.Asset{}

	// Pull asset record from coinbasepro
	//    need to test that asset information is correctly populated
	//    from the coinbase record
	//    currently only defining non FIAT currency as "assets"

	var ledgers []coinbasepro.LedgerEntry

	for _, a := range accounts {
		cursor := client.ListAccountLedger(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledgers); err != nil {
				println(err.Error())
			}
			for _, e := range ledgers {
				if e.Type == "match" {
					currencies := strings.Split(e.Details.ProductID, "-")

					// Determine if this was a buy or sell
					transferAmount, _ := strconv.ParseFloat(e.Amount, 64)
					if transferAmount > float64(0.0) {

						// Store asset
						asset := assetmanagement.Asset{
							ID:       e.Details.TradeID,
							Currency: currencies[0], // Not certain how safe this is
							Quantity: e.Amount,
							BuyDate:  time.Time(e.CreatedAt.Time()),
							BuyPrice: "",
							Cost:     "",
						}
						assets = append(assets, asset)

						// Store purchase price
					} else {
						// Do I need a loop?
						for index, asset := range assets {
							if asset.ID == e.Details.TradeID {
								assets[index].BuyPrice = e.Amount // better way?
							}
						}
					}

					// Store fee
				} else if e.Type == "fee" {
					for index, asset := range assets {
						if asset.ID == e.Details.TradeID {
							assets[index].Cost = e.Amount
						}
					}
				}
			}
		}
	}

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		fmt.Println(message.Price)
	}
}
