package models

type Connection struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	TenantType     string `json:"tenantType"`
	TenantName     string `json:"tenantName"`
	CreateDateUTC  string `json:"createdDateUtc"`
	UpdatedDateUTC string `json:"updatedDateUtc"`
}

type LoginInformation struct {
	Connections  []*Connection `json:"connections"`
	TokenDetails *TokenDetails `json:"token_details"`
}
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}
