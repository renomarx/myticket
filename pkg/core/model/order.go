package model

type Order struct {
	ID       string
	VAT      float64
	Total    float64
	Products []Product
}
