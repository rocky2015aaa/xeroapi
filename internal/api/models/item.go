package models

type Items struct {
	Items []*Item `json:"items"`
}

type Item struct {
	Code         string        `json:"code"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	IsPurchased  bool          `json:"is_purshased"`
	SalesDetails *SalesDetails `json:"sale_details"`
}

type SalesDetails struct {
	UnitPrice   float64 `json:"unit_price"`
	AccountCode string  `json:"account_code"`
}
