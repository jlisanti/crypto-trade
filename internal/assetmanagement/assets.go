package assetmanagement

import (
	"time"
)

// Note that asset quantities are stored in their recorded currency

//type Assets struct {
//	Index  int
//	Assets []asset
//}

type Asset struct {
	ID       string
	Currency string
	Quantity string
	BuyDate  time.Time
	BuyPrice string
	Cost     string
}
