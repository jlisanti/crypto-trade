package main

import (
	//"fmt"
	//"os"
	//"github.com/shopspring/decimal"
	//"github.com/jlisanti/crypto-trade/internal/assetmanagement"
	//"github.com/jlisanti/crypto-trade/internal/finance"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

func main() {

	//var assets = []assetmanagement.Asset{
	//		assetmanagement.Asset{
	//			Index:    0,
	//			Currency: "USD",
	//			Quantity: "100.00",
	//			BuyDate:  "Nov-14-2021",
	//			BuyPrice: "100.00",
	//			Cost:     "0.00",
	//		},
	//	}

	//	println(os.Getenv("COINBASE_PRO_KEY"))
	//	println(os.Getenv("COINBASE_PRO_PASSPHRASE"))
	//	println(os.Getenv("COINBASE_PRO_SECRET"))

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

	for _, a := range accounts {
		println(a.Balance)
	}
}
