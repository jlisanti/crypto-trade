package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"os"

	//"github.com/shopspring/decimal"

	"github.com/Arafatk/glot"

	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	"github.com/jlisanti/crypto-trade/internal/finance"
	"github.com/jlisanti/crypto-trade/pkg/utilities"

	ws "github.com/gorilla/websocket"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

type Points struct {
	x float64
	y float64
}

func main() {

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// configure coinbasepro

	client := coinbasepro.NewClient()

	// Sandbox key
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

	BTCavg := utilities.NewMovingAverage(1.0)

	// Logic for determining slope FIND A BETTER WAY
	slp1 := 0.0
	slp2 := 0.0
	//t1 := time.Now()
	//t3 := time.Now()
	slpIndex := 0
	slpSpacing := 3
	slope := 0.0

	fill1 := true

	// glot settings
	dimensions := 2
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		// Loop across assets and compute the ROI
		//for index, asset := range assets {
		newPrice, _ := strconv.ParseFloat(message.Price, 64)
		newTime := message.Time

		// First newPrice message is always ZERO, need better way to prevent this from going through
		if newPrice != 0.0 {
			utilities.UpdateValue(BTCavg, newPrice, time.Time(newTime))

			if fill1 == true {
				slp1 = BTCavg.AverageValue
				//t1 = time.Time(newTime)
				fill1 = false
			} else if slpIndex == slpSpacing {
				slp2 = BTCavg.AverageValue
				//t3 = time.Time(newTime)
				fill1 = true

				slpIndex = 0
				slope = (slp2 - slp1)
			} else {
				slpIndex++
			}

			//dt := t3.Sub(t1)

			roi, value, age := finance.ComputeROI(message.Price, assets[0].Quantity, assets[0].BuyPrice, assets[0].Cost, assets[0].BuyDate)
			fmt.Fprintln(w, "roi: ", roi, " value: ", value, " age: ", age, " average: ", BTCavg.AverageValue, " slope: ", slope)

			pointGroupName := "Price"
			style := "circle"
			//points := zip(BTCavg.Value, BTCavg.Value)
			//points := [][]float64{{7, 3, 3, 5.6, 5.6, 7, 7, 9, 13, 13, 9, 9}, {10, 10, 4, 4, 5.4, 5.4, 4, 4, 4, 10, 10, 4}}
			points := [][]float64{{7, 3, 13, 5.6, 11.1}, {12, 13, 11, 1, 7}}
			plot.AddPointGroup(pointGroupName, style, points)
			plot.SetTitle("Example Plot")
			plot.SetXLabel("X-axis")
			plot.SetYLabel("Y-axis")

			plot.SavePlot("2.png")
		}

		//}
	}

	fmt.Fprintln(w, "hello")
}
func zip(ts []float64, vs []float64) []Points {
	if len(ts) != len(vs) {
		panic("not same length")
	}

	var res []Points
	for i, t := range ts {
		res = append(res, Points{x: t, y: vs[i]})
	}

	return res
}
