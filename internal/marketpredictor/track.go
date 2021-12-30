package marketpredictor

import (
	"fmt"
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	"github.com/jlisanti/crypto-trade/internal/finance"
	"github.com/jlisanti/crypto-trade/pkg/utilities"
	"github.com/preichenberger/go-coinbasepro/v2"
)

func TrackMarket(assets []assetmanagement.Asset, client *coinbasepro.Client) {
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
			{
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
	//BTCavg := utilities.NewMovingAverage(1.0)

	BTCavg05hr := utilities.NewMovingAverage(0.5)
	BTCavg1hr := utilities.NewMovingAverage(1.0)
	//BTCavg24hr := utilities.NewMovingAverage(24.0)

	// find a better way
	// should consider linear regression to fit last 3-10 data points
	// this would provide a good local slope estimation
	/*
		slp1 := 0.0
		slp2 := 0.0
		slpIndex := 0
		slpSpacing := 3
		slope := 0.0

		fill1 := true
	*/

	m05hr := 0.0
	c05hr := 0.0

	m1hr := 0.0
	c1hr := 0.0

	//m24hr := 0.0
	//c24hr := 0.0

	slopeNeg05hr := false
	//slopeCheck1hr := false
	//slopeCheck24hr := false

	for {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		newPrice, _ := strconv.ParseFloat(message.Price, 64)
		newTime := message.Time

		if newPrice != 0.0 {
			utilities.UpdateValue(BTCavg05hr, newPrice, time.Time(newTime))
			utilities.UpdateValue(BTCavg1hr, newPrice, time.Time(newTime))
			//utilities.UpdateValue(BTCavg24hr, newPrice, time.Time(newTime))

			/*
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
			*/

			/*
				xValues := BTCavg.TimeValues
				y1Values := strings.Fields(strings.Trim(fmt.Sprint(BTCavg.Averages), "[]"))
				y2Values := strings.Fields(strings.Trim(fmt.Sprint(BTCavg.Value), "[]"))
			*/

			/*
				fmt.Println("Fitting averages")
				m, c = utilities.LinearFit(BTCavg.Averages, BTCavg.TimeDiffs, m, c)
			*/

			m05hr, c05hr = utilities.LinearFit(BTCavg05hr.Averages, BTCavg05hr.TimeDiffs, m05hr, c05hr)
			m1hr, c1hr = utilities.LinearFit(BTCavg1hr.Averages, BTCavg1hr.TimeDiffs, m1hr, c1hr)
			//m24hr, c24hr = utilities.LinearFit(BTCavg24hr.Averages, BTCavg24hr.TimeDiffs, m24hr, c24hr)

			roi, value, age := finance.ComputeROI(message.Price, assets[0].Quantity, assets[0].BuyPrice, assets[0].Cost, assets[0].BuyDate)

			// Prim sell logic
			if roi > 1.0 {
				if m05hr < 0.0 {
					if !slopeNeg05hr {
						slopeNeg05hr = true
						if m1hr < 0.0 {
							//if m24hr < 0.0 {
							// make puchase
							fmt.Println("Making puchase... roi: ", roi, " age: ", age, " value: ", value)
							//}
						}
					}
				}
			}
			accounts, err := client.GetAccounts()
			if err != nil {
				println(err.Error())
			}

			purchaseAmount := 0.0
			for _, a := range accounts {
				if a.Currency == "USD" {
					purchaseAmount, _ = strconv.ParseFloat(a.Balance, 64)
				}
			}

			// Prim buy logic
			if purchaseAmount > 0.0 {
				if m05hr > 0.0 {
					if slopeNeg05hr {
						slopeNeg05hr = false
						if m1hr > 0.0 {
							//if m24hr > 0.0 {
							fmt.Println("Buying bitcoin... ", slopeNeg05hr)
							//}
						}
					}
				}

			}
			if m05hr > 0.0 {
				slopeNeg05hr = false
			}

			fmt.Println("price: ", message.Price, " roi: ", roi, " value: ", value, "age: ", age, "average: ", BTCavg05hr.AverageValue, " m05hr: ", m05hr, " m1hr: ", m1hr) //, " m24hr: ", m24hr)

			/*
				//defer file.Close()

					file, err := os.Create("./dat/data.csv")
					if err != nil {
						log.Fatalln("Failed to open file", err)
					}

					w := csv.NewWriter(file)
					//defer w.Flush()

					for i, xValue := range xValues {
						//row := []string{xValue, y1Values[i], y2Values[i]}
						timeDifference := fmt.Sprintf("%f", BTCavg.TimeDiffs[i])
						row := []string{xValue, y1Values[i], y2Values[i], timeDifference}
						if err := w.Write(row); err != nil {
							log.Fatalln("error writing record to file", err)
						}
					}
					w.Flush()
					file.Close()
			*/
		}
		time.Sleep(1e9)
	}
}
