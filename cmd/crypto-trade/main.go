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
	router.ServeFiles("/dat/*filepath", http.Dir("dat"))

	router.GET("/", serveHomepage)
	client := coinbasepro.ConnectCoinbasepro(&assets)
	fmt.Println(assets[0].BuyPrice)

	accounts, err := client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	for _, a := range accounts {
		println(a.Currency, a.Balance)
	}

	go marketpredictor.TrackMarket(assets, &client)

	err = http.ListenAndServe(":8080", router)
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
