package assetmanagement

import ()

// Note that asset quantities are stored in their recorded currency

//type Assets struct {
//	Index  int
//	Assets []asset
//}

type Asset struct {
	Index    int
	Currency string
	Quantity string
	BuyDate  string
	BuyPrice string
	Cost     string
}
