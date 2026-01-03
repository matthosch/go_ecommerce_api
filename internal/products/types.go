package products

type Product struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Quantity     int32  `json:"quantity"`
}
