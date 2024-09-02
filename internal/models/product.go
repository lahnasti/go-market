package models

type Product struct {
	UID int `json:"uid"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price float64 `json:"price"`
}