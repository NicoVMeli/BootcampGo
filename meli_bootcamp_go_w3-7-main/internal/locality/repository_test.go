package locality

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateOk(t *testing.T) {
	localityToSave := domain.Locality{
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("INSERT INTO localities").
		ExpectExec().
		WithArgs(localityToSave.ZipCode, localityToSave.LocalityName, localityToSave.ProvinceName, localityToSave.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	localityRepository := NewRepository(db)
	actualId, err := localityRepository.Save(context.Background(), localityToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateConflict(t *testing.T) {
	localityToSave := domain.Locality{
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("INSERT INTO localities").
		ExpectExec().
		WithArgs(localityToSave.ZipCode, localityToSave.LocalityName, localityToSave.ProvinceName, localityToSave.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1))
	localityRepository := NewRepository(db)
	ctx := context.Background()
	actualId, err := localityRepository.Save(ctx, localityToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	localityCompare := domain.Locality{
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}
	mock.
		ExpectPrepare("INSERT INTO localities").
		ExpectExec().
		WithArgs(localityToSave.ZipCode, localityCompare.LocalityName, localityCompare.ProvinceName, localityCompare.CountryName).
		WillReturnError(errors.New("zipCode duplicated"))
	ctx = context.Background()
	actualId, err = localityRepository.Save(ctx, localityCompare)
	assert.NotNil(t, err)
	assert.Equal(t, 0, actualId)
	assert.Equal(t, localityToSave.ZipCode, localityCompare.ZipCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}
