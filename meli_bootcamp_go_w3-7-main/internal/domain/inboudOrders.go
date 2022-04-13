package domain

import "time"

type InboudOrders struct {
	ID             int       `json:"id"`
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeId     int       `json:"employee_id"`
	ProductBatchId int       `json:"product_batch_id"`
	WarehouseID    int       `json:"warehouse_id"`
}
