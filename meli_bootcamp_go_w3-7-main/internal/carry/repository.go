package carry

import (
	"context"
	"database/sql"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Repository encapsulates the storage of a carry.
type Repository interface {
	Close()
	GetAll(ctx context.Context) ([]domain.Carry, error)
	Get(ctx context.Context, id int) (domain.Carry, error)
	Exists(ctx context.Context, carryId string) bool
	Save(ctx context.Context, c domain.Carry) (int, error)
	Update(ctx context.Context, c domain.Carry) error
	Delete(ctx context.Context, id int) error
	GetCarriesByLocality(ctx context.Context) ([]domain.CarriesReport, error)
	GetCarriesByLocalityId(ctx context.Context, id int) (domain.CarriesReport, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Close() {
	r.db.Close()
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Carry, error) {
	query := "SELECT * FROM carries"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var carries []domain.Carry

	for rows.Next() {
		c := domain.Carry{}
		_ = rows.Scan(&c.CID, &c.CompanyName, &c.Address, &c.Telephone, &c.BatchNumber, &c.LocalityId, &c.ID)
		carries = append(carries, c)
	}

	return carries, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Carry, error) {
	query := "SELECT * FROM carries WHERE id=?;"
	row := r.db.QueryRow(query, id)
	c := domain.Carry{}
	err := row.Scan(&c.CID, &c.CompanyName, &c.Address, &c.Telephone, &c.BatchNumber, &c.LocalityId, &c.ID)
	if err != nil {
		return domain.Carry{}, err
	}

	return c, nil
}

func (r *repository) Exists(ctx context.Context, carryId string) bool {
	query := "SELECT cid FROM carries WHERE cid=?;"
	row := r.db.QueryRow(query, carryId)
	err := row.Scan(&carryId)
	return err == nil
}

func (r *repository) Save(ctx context.Context, c domain.Carry) (int, error) {
	query := "INSERT INTO carries (cid, company_name, address, telephone, locality_id, batch_number) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&c.CID, &c.CompanyName, &c.Address, &c.Telephone, &c.BatchNumber, &c.LocalityId)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, c domain.Carry) error {
	query := "UPDATE carries SET cid=?, company_name=?, address=?, telephone=?, batch_number=?, locality_id=? WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&c.CID, &c.CompanyName, &c.Address, &c.Telephone, &c.BatchNumber, &c.LocalityId, &c.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM carries WHERE id=?"
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

func (r *repository) GetCarriesByLocality(ctx context.Context) ([]domain.CarriesReport, error) {
	query := "SELECT locality_id, locality_name, COUNT(carries.id) AS carries_count FROM carries INNER JOIN localities ON carries.locality_id = localities.id GROUP BY locality_id, locality_name;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var report []domain.CarriesReport

	for rows.Next() {
		c := domain.CarriesReport{}
		_ = rows.Scan(&c.LocalityId, &c.LocalityName, &c.CarriesCount)
		report = append(report, c)
	}
	return report, nil
}

func (r *repository) GetCarriesByLocalityId(ctx context.Context, id int) (domain.CarriesReport, error) {
	query := "SELECT locality_id, locality_name, COUNT(carries.id) AS carries_count FROM carries INNER JOIN localities ON locality_id = localities.id WHERE locality_id = ? GROUP BY locality_id, locality_name;"
	row := r.db.QueryRow(query, id)
	c := domain.CarriesReport{}
	err := row.Scan(&c.LocalityId, &c.LocalityName, &c.CarriesCount)
	if err != nil {
		return domain.CarriesReport{}, err
	}

	return c, nil
}
