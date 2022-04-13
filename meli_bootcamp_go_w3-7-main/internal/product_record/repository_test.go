package product_record

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProductRecord_CreateOk(t *testing.T) {
	// ARRANGE ----------------------------------------
	LocalTimezoneDate := time.Now()            // Time is created using local system timezone (-03)
	UTCTimezoneDate := LocalTimezoneDate.UTC() // Time is converted to standard UTC timezone
	productRecordToSave := domain.ProductRecord{
		ID:             0,
		LastUpdateDate: &UTCTimezoneDate,
		SalePrice:      23.1,
		PurchasePrice:  23.5,
		ProductID:      1,
	}

	// SQL Mock initialization
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	// Set expected queries and respective returns
	mock.
		ExpectPrepare("INSERT INTO product_records").
		ExpectExec().
		WithArgs(productRecordToSave.LastUpdateDate, 23.5, 23.1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	// Repository initialization
	productRecordRepository := NewRepository(db)

	// ACT ----------------------------------------
	actualId, err := productRecordRepository.Save(context.Background(), productRecordToSave)

	//ASSERT --------------------------------------
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRecord_GetOk(t *testing.T) {
	// ARRANGE ----------------------------------------
	date := time.Date(2022, 02, 07, 19, 57, 52, 0, time.UTC) // Time is created using standard UTC timezone, as expected to be returned from the repository
	expectedProductRecord := domain.ProductRecord{
		ID:             5,
		LastUpdateDate: &date,
		SalePrice:      23.1,
		PurchasePrice:  23.5,
		ProductID:      1,
	}

	// SQL Mock initialization
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	// Set mock columns
	columns := []string{"id", "last_update_date", "purchase_price", "sale_price", "product_id"}
	rows := sqlmock.NewRows(columns)
	// Set Product Record id to search
	productRecordId := 5
	// Add row to the mock. The datetime is in UTC timezone
	rows.AddRow(productRecordId, `2022-02-07 19:57:52`, 23.5, 23.1, 1)
	// Set expected query and respective returns
	mock.
		ExpectPrepare("SELECT id, last_update_date, purchase_price, sale_price, product_id FROM product_records").
		ExpectQuery().
		WithArgs(5).
		WillReturnRows(rows)
	// Repository initialization
	productRecordRepository := NewRepository(db)

	// ACT ----------------------------------------
	actualResult, err := productRecordRepository.Get(context.Background(), productRecordId)

	//ASSERT --------------------------------------
	fmt.Println(actualResult)
	assert.Nil(t, err)
	assert.NotNil(t, actualResult)
	assert.Equal(t, expectedProductRecord, actualResult)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRecord_CreateConflict(t *testing.T) {
	// ARRANGE ----------------------------------------
	futureDateUTC := time.Date(4000, 02, 07, 19, 57, 52, 0, time.UTC)
	productRecordToSave := domain.ProductRecord{
		ID:             0,
		LastUpdateDate: &futureDateUTC,
		SalePrice:      23.1,
		PurchasePrice:  23.5,
		ProductID:      1,
	}

	// SQL Mock initialization
	db, _, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	// Repository initialization
	productRecordRepository := NewRepository(db)

	// ACT ----------------------------------------
	actualId, err := productRecordRepository.Save(context.Background(), productRecordToSave)

	//ASSERT --------------------------------------
	assert.NotNil(t, err)
	assert.Equal(t, 0, actualId)
	assert.ErrorIs(t, err, ErrLastUpdateDateExceedsCurrentDate)
}
