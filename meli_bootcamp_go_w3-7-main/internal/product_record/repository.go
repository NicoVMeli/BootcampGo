package product_record

import (
	"context"
	"database/sql"
	"errors"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id int) (domain.ProductRecord, error)
	Save(ctx context.Context, record domain.ProductRecord) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// Errors
var (
	ErrLastUpdateDateExceedsCurrentDate = errors.New("last update date provided exceeds the current system date")
)

// DB queries
const (
	saveQuery = "INSERT INTO product_records(last_update_date, purchase_price, sale_price, product_id)" +
		"VALUES (?, ?, ?, ?)"
	getQuery = "SELECT id, last_update_date, purchase_price, sale_price, product_id " +
		"FROM product_records pr " +
		"WHERE pr.id=?"
)

func (r *repository) Save(ctx context.Context, productRecord domain.ProductRecord) (int, error) {
	// Validate that productRecord.last_update_date < current system date
	currentSystemDate := time.Now()
	isAfter := productRecord.LastUpdateDate.After(currentSystemDate)
	if isAfter {
		return 0, ErrLastUpdateDateExceedsCurrentDate
	}
	// Prepare query
	stmt, err := r.db.Prepare(saveQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	// productRecord.LastUpdateDate is automatically converted to UTC timezone.
	// When you are saving the value, the mysql driver it first converts the time stamp to the UTC time zone and then sends it off to the database.
	result, err := stmt.Exec(productRecord.LastUpdateDate, productRecord.PurchasePrice, productRecord.SalePrice, productRecord.ProductID)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertId), nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.ProductRecord, error) {
	stmt, err := r.db.Prepare(getQuery)
	if err != nil {
		return domain.ProductRecord{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err != nil {
		return domain.ProductRecord{}, err
	}
	productRecord := domain.ProductRecord{}
	var lastUpdateDateInStringFormat string
	// Scan last_update_date to a variable of type string, and scan the rest of the data to the respective ProductRecord object fields
	err = row.Scan(&productRecord.ID, &lastUpdateDateInStringFormat, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)
	if err != nil {
		return domain.ProductRecord{}, err
	}
	// Parse last_update_time from string format to Time format (UTC)
	lastUpdateDateInTimeFormatUTC, err := time.Parse("2006-01-02 15:04:05", lastUpdateDateInStringFormat)
	// Assign last_update_time (UTC) to respective ProductRecord object field
	productRecord.LastUpdateDate = &lastUpdateDateInTimeFormatUTC
	return productRecord, nil
}
