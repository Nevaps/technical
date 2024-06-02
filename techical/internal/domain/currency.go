package domain

type CurrencyResponse struct {
	Result float64 `json:"result"`
}

type Currency struct {
	Name      string `json:"name"`
	Ticker    string `json:"ticker"`
	Available bool   `json:"available"`
}
