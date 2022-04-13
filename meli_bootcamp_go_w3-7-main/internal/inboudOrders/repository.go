package inboudOrders

import (
	"context"
	"database/sql"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, i domain.InboudOrders) (int, error)
	ExistsEmployee(ctx context.Context, employeeId int) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, i domain.InboudOrders) (int, error) {
	query := "INSERT INTO inbound_orders(order_date,order_number,employee_id,product_batch_id,warehouse_id) VALUES (?,?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&i.OrderDate, &i.OrderNumber, &i.EmployeeId, &i.ProductBatchId, &i.WarehouseID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) ExistsEmployee(ctx context.Context, employeeId int) bool {
	query := "SELECT id FROM employees WHERE id=?;"
	row := r.db.QueryRow(query, employeeId)
	err := row.Scan(&employeeId)
	return err == nil
}
