package productBatches

import (
	"context"
	"database/sql"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Repository encapsulates the storage of a section.
type Repository interface {
	GetBySectionNumber(ctx context.Context, sectionNumber int) (domain.Section, error)
	GetQuantity(ctx context.Context, quant int) (domain.ProductBatches, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.ProductBatches) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, s domain.ProductBatches) (int, error) {
	query := "INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&s.BatchNumber, &s.CurrentQuantity, &s.CurrentTemperature, &s.DueDate, &s.InitialQuantity, &s.ManufacturingDate, &s.ManufacturingHour, &s.MinimumTemperature, &s.ProductId, &s.SectionId)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Exists(ctx context.Context, batchNumber int) bool {
	query := "SELECT batch_number FROM product_batches WHERE batch_number=?;"
	row := r.db.QueryRow(query, batchNumber)
	err := row.Scan(&batchNumber)
	return err == nil
}

func (r *repository) GetBySectionNumber(ctx context.Context, sectionNumber int) (domain.Section, error) {
	query := "SELECT * FROM sections WHERE section_number=?"
	row := r.db.QueryRow(query, sectionNumber)
	s := domain.Section{}
	err := row.Scan(&s.SectionNumber)
	if err != nil {
		return domain.Section{}, err
	}
	return s, nil
}

func (r *repository) GetQuantity(ctx context.Context, quant int) (domain.ProductBatches, error) {
	query := "SELECT * FROM product_batches WHERE current_quantity=?"
	row := r.db.QueryRow(query, quant)
	s := domain.ProductBatches{}
	err := row.Scan(&s.CurrentQuantity, &s.ProductId)
	if err != nil {
		return domain.ProductBatches{}, err
	}
	return s, nil
}
