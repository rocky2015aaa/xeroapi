package models

type Contacts struct {
	Items []*Item `json:"contacts"`
}

type Contact struct {
	ContactStatus      string        `json:"contract_status"`
	Name               string        `json:"name"`
	EmailAddress       string        `json:"email_address"`
	BankAccountDetails string        `json:"bank_account_details"`
	Addresses          []*Address    `json:"addresses"`
	Phones             []*Phone      `json:"phones"`
	UpdatedDateUTC     string        `json:"update_date_utc"`
	IsSupplier         bool          `json:"is_supplier"`
	IsCustomer         bool          `json:"is_customer"`
	PaymentTerms       *PaymentTerms `json:""`
}

type Address struct {
	AddressType string `json:"address_type"`
	City        string `json:"city"`
	Region      string `json:"region"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

type Phone struct {
	PhoneType        string `json:"phone_type"`
	PhoneNumber      string `json:"phone_number"`
	PhoneAreaCode    string `json:"phone_area_code"`
	PhoneCountryCode string `json:"phone_country_code"`
}

type PaymentTerms struct {
	Bills *PaymentTermsData `json:"bills"`
	Sales *PaymentTermsData `json:"sales"`
}

type PaymentTermsData struct {
	Day  int    `json:"day"`
	Type string `json:"type"`
}
