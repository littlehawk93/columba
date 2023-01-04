package fedex

type propertiesResponse struct {
	API propertiesAPI `json:"api"`
}

type propertiesAPI struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
