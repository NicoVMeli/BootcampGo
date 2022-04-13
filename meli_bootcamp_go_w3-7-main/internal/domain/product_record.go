package domain

import "time"

type ProductRecord struct {
	ID             int        `json:"id"`
	LastUpdateDate *time.Time `json:"last_update_date" binding:"required"`
	PurchasePrice  float64    `json:"purchase_price" binding:"required"`
	SalePrice      float64    `json:"sale_price" binding:"required"`
	ProductID      int        `json:"product_id" binding:"required"`
}
