package carry

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	r "github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
)

var c = &r.Carry{
	ID:          1,
	CID:         "c1",
	CompanyName: "company1",
	Address:     "Address1",
	Telephone:   "123456",
	LocalityId:  1,
	BatchNumber: 1,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestCreateCarryOK(t *testing.T) {
	db, mock := NewMock()
	repo := &repository{db}
	defer func() {
		repo.Close()
	}()
	carry := domain.Carry{
		ID:          1,
		CID:         "c1",
		CompanyName: "company1",
		Address:     "Address1",
		Telephone:   "123456",
		LocalityId:  1,
		BatchNumber: 1,
	}
	query := "INSERT INTO carries"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(c.CID, c.CompanyName, c.Address, c.Telephone, c.LocalityId, c.BatchNumber).WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := repo.Save(context.Background(), carry)
	assert.NoError(t, err)
}

func TestCreateCarryConflict(t *testing.T) {
	expectedMessageError := "El cid ingresado ya existe en la base de datos"
	db, mock := NewMock()
	repo := &repository{db}
	defer func() {
		repo.Close()
	}()
	carry := domain.Carry{
		ID:          1,
		CID:         "c1",
		CompanyName: "company1",
		Address:     "Address1",
		Telephone:   "123456",
		LocalityId:  1,
		BatchNumber: 1,
	}
	query := "INSERT INTO carries"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(c.CID, c.CompanyName, c.Address, c.Telephone, c.LocalityId, c.BatchNumber).WillReturnError(errors.New(expectedMessageError))

	_, err := repo.Save(context.Background(), carry)

	assert.Equal(t, expectedMessageError, err.Error())
}
