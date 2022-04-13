package domain

// struc ProductBatches req 3
type ProductBatches struct {
	ID                 int     `json:"id"`
	BatchNumber        int     `json:"batch_number"`
	CurrentQuantity    int     `json:"current_quantity"`
	CurrentTemperature int     `json:"current_temperature"`
	DueDate            string  `json:"due_date"`
	InitialQuantity    int     `json:"initial_quantity"`
	ManufacturingDate  string  `json:"manufacturing_date"`
	ManufacturingHour  int     `json:"manufacturing_hour"`
	MinimumTemperature float32 `json:"minimum_temperature"`
	ProductId          int     `json:"product_id"`
	SectionId          int     `json:"section_id"`
}
