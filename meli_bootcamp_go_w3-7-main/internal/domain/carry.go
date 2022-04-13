package domain

type Carry struct {
	ID          int    `json:"id"`
	CID         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
	BatchNumber int    `json:"batch_number"`
}

type CarriesReport struct {
	LocalityId   int    `jsn:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int    `json:"carries_count"`
}

type NewCarry struct {
	ID          int    `json:"id" binding:"required"`
	CID         string `json:"cid" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Telephone   string `json:"telephone" binding:"required"`
	LocalityId  int    `json:"locality_id" binding:"required"`
	BatchNumber int    `json:"batch_number" binding:"required"`
}
