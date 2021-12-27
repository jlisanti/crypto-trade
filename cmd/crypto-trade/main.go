package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	"github.com/jlisanti/crypto-trade/internal/coinbasepro"
	"github.com/jlisanti/crypto-trade/internal/marketpredictor"
	"github.com/julienschmidt/httprouter"

	//"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
)

type HomePage struct {
	Time string
}

const serviceName = "CyptoTrade-0.01"
const serviceDescription = "Crypto trade bot"

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
	assets           = []assetmanagement.Asset{}
)

func serveHomepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writingSync.Lock()
	programIsRunning = true
	writingSync.Unlock()

	var homepage HomePage
	homepage.Time = time.Now().Format("02/01/2006, 15:03:04")

	tmpl := template.Must(template.ParseFiles("html/index.html"))
	_ = tmpl.Execute(writer, homepage)

	writingSync.Lock()
	programIsRunning = false
	writingSync.Unlock()
}

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()
	for programIsRunning {
		fmt.Println(s.String() + " stopping...")
	}
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	router := httprouter.New()
	//timer := sse.New()

	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.GET("/", serveHomepage)
	coinbasepro.ConnectCoinbasepro(&assets)
	fmt.Println(assets[0].BuyPrice)

	go marketpredictor.TrackMarket(assets)

	/*
		router.POST("/get_time", getTime)

		router.Handler("GET", "/time", timer)
		go streamTime(timer)
	*/

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Problem starting web server: " + err.Error())
		os.Exit(-1)
	}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}

	prg := &program{}
	s, err := service.New(prg, serviceConfig)

	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}

	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}
}

/*
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

	for {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		// Loop across assets and compute the ROI
		//for index, asset := range assets {
		//newPrice, _ := strconv.ParseFloat(message.Price, 64)
		//newTime := message.Time

		// First newPrice message is always ZERO, need better way to prevent this from going through
		if newPrice != 0.0 {
			utilities.UpdateValue(BTCavg, newPrice, time.Time(newTime))

			if fill1 {
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

			fmt.Fprintf(w, `<html>
			<head>
			<script type="text/javascript"
			  src="/Users/Joel/sources/dygraphs/dygraph.js"></script>
			<link rel="stylesheet" src="dygraph.css" />
			</head>
			<body>
			<div id="graphdiv2"
			  style="width:500px; height:300px;"></div>
			<script type="text/javascript">
			  g2 = new Dygraph(
				document.getElementById("graphdiv2"),
				"/Users/Joel/projects/crypto-trade/temperatures.csv", // path to CSV file
				{}          // options
			  );
			</script>
			</body>
			</html>`)
			/*
							hello := "hello yu bn"
							fmt.Fprintf(w, `<html>
				            <head>
				            </head>
				            <body>
				            <h1>Go Timer (ticks every second!)</h1>
				            <div id="output"></div>
				            <script type="text/javascript">
				            console.log("`+hello+`");
				            </script>
				            </body>

			//roi, value, age := finance.ComputeROI(message.Price, assets[0].Quantity, assets[0].BuyPrice, assets[0].Cost, assets[0].BuyDate)
			//fmt.Fprintln(w, "roi: ", roi, " value: ", value, " age: ", age, " average: ", BTCavg.AverageValue, " slope: ", slope)
		}
	}
}
*/
