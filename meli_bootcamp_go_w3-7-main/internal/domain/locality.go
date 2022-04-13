package domain

type Locality struct {
	ID           int    `json:"id"`
	ZipCode      string `json:"zip_code"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}
