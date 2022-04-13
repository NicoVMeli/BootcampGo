package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Repository encapsulates the storage of a employee.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Get(ctx context.Context, id int) (domain.Employee, error)
	Exists(ctx context.Context, cardNumberID string) bool
	Save(ctx context.Context, e domain.Employee) (int, error)
	Update(ctx context.Context, e domain.Employee) error
	Delete(ctx context.Context, id int) error
	GetAllReportInboundOrders() ([]domain.ReportInboundOrders, error)
	GetReportByEmployeeIdInboundOrders(employeeId int) (domain.ReportInboundOrders, error)
}

const QueryReportInboundOrders = `select e.id, e.card_number_id,e.first_name,e.last_name,e.warehouse_id,count(i.id) as inbound_orders_count from employees e
inner join inbound_orders i on e.id = i.employee_id
 %s 
group by e.id;`

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	query := "SELECT * FROM employees"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var employees []domain.Employee

	for rows.Next() {
		e := domain.Employee{}
		_ = rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
		employees = append(employees, e)
	}

	return employees, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Employee, error) {
	query := "SELECT * FROM employees WHERE id=?;"
	row := r.db.QueryRow(query, id)
	e := domain.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if err != nil {
		return domain.Employee{}, err
	}

	return e, nil
}

func (r *repository) Exists(ctx context.Context, cardNumberID string) bool {
	query := "SELECT card_number_id FROM employees WHERE card_number_id=?;"
	row := r.db.QueryRow(query, cardNumberID)
	err := row.Scan(&cardNumberID)
	return err == nil
}

func (r *repository) Save(ctx context.Context, e domain.Employee) (int, error) {
	query := "INSERT INTO employees(card_number_id,first_name,last_name,warehouse_id) VALUES (?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, e domain.Employee) error {
	query := "UPDATE employees SET first_name=?, last_name=?, warehouse_id=?  WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&e.FirstName, &e.LastName, &e.WarehouseID, &e.ID)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM employees WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) GetAllReportInboundOrders() ([]domain.ReportInboundOrders, error) {
	query := fmt.Sprintf(QueryReportInboundOrders, "")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var repostEmployees []domain.ReportInboundOrders

	for rows.Next() {
		i := domain.ReportInboundOrders{}
		_ = rows.Scan(&i.ID, &i.CardNumberID, &i.FirstName, &i.LastName, &i.WarehouseID, &i.InboundOrdersCount)
		repostEmployees = append(repostEmployees, i)
	}

	return repostEmployees, nil
}

func (r *repository) GetReportByEmployeeIdInboundOrders(employeeId int) (domain.ReportInboundOrders, error) {
	query := fmt.Sprintf(QueryReportInboundOrders, "where e.id = ?")
	row := r.db.QueryRow(query, employeeId)

	i := domain.ReportInboundOrders{}
	err := row.Scan(&i.ID, &i.CardNumberID, &i.FirstName, &i.LastName, &i.WarehouseID, &i.InboundOrdersCount)
	if err != nil {
		return domain.ReportInboundOrders{}, err
	}

	return i, nil
}
