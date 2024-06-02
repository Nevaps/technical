package structs

type FiatPrices struct {
	Base   string             `json:"base"`
	Prices map[string]float64 `json:"results"`
}

type CryptoPrices struct {
	Prices map[string]float64 `json:"prices"`
}
