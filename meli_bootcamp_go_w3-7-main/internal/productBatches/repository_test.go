package productBatches

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateOkProductBatches(t *testing.T) {
	ProdBatchesToSave := domain.ProductBatches{
		BatchNumber:        111,
		CurrentQuantity:    200,
		CurrentTemperature: 20,
		DueDate:            "2022-04-04",
		InitialQuantity:    10,
		ManufacturingDate:  "2020-04-04",
		ManufacturingHour:  10,
		MinimumTemperature: 5,
		ProductId:          1,
		SectionId:          1,
	}
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("INSERT INTO product_batches").
		ExpectExec().
		WithArgs(ProdBatchesToSave.BatchNumber,
			ProdBatchesToSave.CurrentQuantity,
			ProdBatchesToSave.CurrentTemperature,
			ProdBatchesToSave.DueDate,
			ProdBatchesToSave.InitialQuantity,
			ProdBatchesToSave.ManufacturingDate,
			ProdBatchesToSave.ManufacturingHour,
			ProdBatchesToSave.MinimumTemperature,
			ProdBatchesToSave.ProductId,
			ProdBatchesToSave.SectionId).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	productBatchesRepository := NewRepository(db)
	actualId, err := productBatchesRepository.Save(context.Background(), ProdBatchesToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateConflictProductBatches(t *testing.T) {
	ProdBatchesToSave := domain.ProductBatches{
		BatchNumber:        111,
		CurrentQuantity:    200,
		CurrentTemperature: 20,
		DueDate:            "2022-04-04",
		InitialQuantity:    10,
		ManufacturingDate:  "2020-04-04",
		ManufacturingHour:  10,
		MinimumTemperature: 5,
		ProductId:          1,
		SectionId:          1,
	}
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("INSERT INTO product_batches").
		ExpectExec().
		WithArgs(ProdBatchesToSave.BatchNumber,
			ProdBatchesToSave.CurrentQuantity,
			ProdBatchesToSave.CurrentTemperature,
			ProdBatchesToSave.DueDate,
			ProdBatchesToSave.InitialQuantity,
			ProdBatchesToSave.ManufacturingDate,
			ProdBatchesToSave.ManufacturingHour,
			ProdBatchesToSave.MinimumTemperature,
			ProdBatchesToSave.ProductId,
			ProdBatchesToSave.SectionId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	productBatchesRepository := NewRepository(db)
	ctx := context.Background()
	actualId, err := productBatchesRepository.Save(ctx, ProdBatchesToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	prodBatchesCompare := domain.ProductBatches{
		BatchNumber:        111,
		CurrentQuantity:    200,
		CurrentTemperature: 20,
		DueDate:            "2022-04-04",
		InitialQuantity:    10,
		ManufacturingDate:  "2020-04-04",
		ManufacturingHour:  10,
		MinimumTemperature: 5,
		ProductId:          1,
		SectionId:          1,
	}
	mock.
		ExpectPrepare("INSERT INTO product_batches").
		ExpectExec().
		WithArgs(prodBatchesCompare.BatchNumber,
			prodBatchesCompare.CurrentQuantity,
			prodBatchesCompare.CurrentTemperature,
			prodBatchesCompare.DueDate,
			prodBatchesCompare.InitialQuantity,
			prodBatchesCompare.ManufacturingDate,
			prodBatchesCompare.ManufacturingHour,
			prodBatchesCompare.MinimumTemperature,
			prodBatchesCompare.ProductId,
			prodBatchesCompare.SectionId).
		WillReturnError(errors.New("batch Number duplicated"))
	ctx = context.Background()
	actualId, err = productBatchesRepository.Save(ctx, prodBatchesCompare)
	assert.NotNil(t, err)
	assert.Equal(t, 0, actualId)
	assert.Equal(t, ProdBatchesToSave.BatchNumber, prodBatchesCompare.BatchNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}
