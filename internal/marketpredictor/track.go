package marketpredictor

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	"github.com/jlisanti/crypto-trade/internal/finance"
	"github.com/jlisanti/crypto-trade/pkg/utilities"
	"github.com/preichenberger/go-coinbasepro/v2"
)

func datePlusTime(date, timeOfDay string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05.000", date+" "+timeOfDay)
}

func TrackMarket(assets []assetmanagement.Asset) {
	fmt.Println("starting track loop")

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	fmt.Println("subscrubing to message channel")

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

	// find a better way
	slp1 := 0.0
	slp2 := 0.0
	slpIndex := 0
	slpSpacing := 3
	slope := 0.0

	fill1 := true

	for {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		newPrice, _ := strconv.ParseFloat(message.Price, 64)
		newTime := message.Time

		if newPrice != 0.0 {
			utilities.UpdateValue(BTCavg, newPrice, time.Time(newTime))

			if fill1 {
				slp1 = BTCavg.AverageValue
				fill1 = false
			} else if slpIndex == slpSpacing {
				slp2 = BTCavg.AverageValue
				fill1 = true

				slpIndex = 0
				slope = (slp2 - slp1)
			} else {
				slpIndex++
			}

			roi, value, age := finance.ComputeROI(message.Price, assets[0].Quantity, assets[0].BuyPrice, assets[0].Cost, assets[0].BuyDate)
			fmt.Println("price: ", message.Price, " roi: ", roi, " value: ", value, "age: ", age, "average: ", BTCavg.AverageValue, " slope: ", slope)

			xValues := BTCavg.TimeValues
			y1Values := strings.Fields(strings.Trim(fmt.Sprint(BTCavg.Averages), "[]"))
			y2Values := strings.Fields(strings.Trim(fmt.Sprint(BTCavg.Value), "[]"))

			fmt.Println("lengths: ", len(xValues), len(y1Values), len(y2Values))
			//defer file.Close()

			file, err := os.Create("./dat/data.csv")
			if err != nil {
				log.Fatalln("Failed to open file", err)
			}

			w := csv.NewWriter(file)
			//defer w.Flush()

			for i, xValue := range xValues {
				//row := []string{xValue, y1Values[i], y2Values[i]}
				row := []string{xValue, y1Values[i], y2Values[i]}
				fmt.Println(row)
				if err := w.Write(row); err != nil {
					log.Fatalln("error writing record to file", err)
				}
			}
			w.Flush()
			file.Close()
		}
		time.Sleep(1e9)
	}
}
