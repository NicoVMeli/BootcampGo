package section

import (
	"context"
	"errors"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoM struct {
	mock.Mock
}

func (r *repoM) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}
func (r *repoM) Get(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}
func (r *repoM) Exists(ctx context.Context, sectionNumber int) bool {
	args := r.Called(ctx, sectionNumber)
	return args.Bool(0)
}
func (r *repoM) Save(ctx context.Context, se domain.Section) (int, error) {
	args := r.Called(ctx, se)
	return args.Int(0), args.Error(1)
}
func (r *repoM) Update(ctx context.Context, se domain.Section) error {
	args := r.Called(ctx, se)
	return args.Error(0)
}
func (r *repoM) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func TestCreateOk(t *testing.T) {
	repo := new(repoM)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)
	s := NewService(repo)
	objetoAPersistir := domain.Section{
		SectionNumber:      41,
		CurrentTemperature: 3,
		MinimumTemperature: 3,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    3,
		WarehouseID:        3,
		ProductTypeID:      3,
	}
	ns, err := s.Save(context.Background(), objetoAPersistir)
	assert.NoError(t, err)
	assert.Equal(t, 1, ns.ID)
	assert.Equal(t, 41, ns.SectionNumber)
}

func TestCreateConflict(t *testing.T) {
	expectedError := errors.New("la seccion ya existe")
	expectedResult := domain.Section{}
	repo := new(repoM)
	repo.On("Exists", mock.Anything, mock.Anything).Return(true)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, err := serviceT.Save(ctx, expectedResult)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, resul)
}

func TestFindAll(t *testing.T) {
	repo := new(repoM)
	repo.On("GetAll", mock.Anything).Return([]domain.Section{
		{
			SectionNumber:      41,
			CurrentTemperature: 3,
			MinimumTemperature: 3,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    3,
			WarehouseID:        3,
			ProductTypeID:      3,
		},
		{
			SectionNumber:      42,
			CurrentTemperature: 3,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		},
	}, nil)
	s := NewService(repo)
	objetosRecuperados, err := s.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, objetosRecuperados, 2)
}

func TestFindOne(t *testing.T) {
	expectedResult := domain.Section{
		ID:                 1,
		SectionNumber:      41,
		CurrentTemperature: 3,
		MinimumTemperature: 3,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    3,
		WarehouseID:        3,
		ProductTypeID:      3,
	}
	repo := new(repoM)
	repo.On("Get", mock.Anything, 1).Return(expectedResult, nil)
	s := NewService(repo)
	ctx := context.Background()
	resul, _ := s.Get(ctx, 1)
	assert.Equal(t, expectedResult, resul)
}

func TestFindOneNoExist(t *testing.T) {
	repo := new(repoM)
	repo.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, errors.New("No existe un section con ese id"))
	s := NewService(repo)
	objetoRecuperado, err := s.Get(context.Background(), 3)
	obj := domain.Section{}
	assert.Error(t, err)
	assert.Equal(t, obj, objetoRecuperado, "No existe un section con ese id")
}

func TestUpdateOk(t *testing.T) {
	repo := new(repoM)
	objetoAc := domain.Section{
		SectionNumber:      41,
		CurrentTemperature: 3,
		MinimumTemperature: 4,
		CurrentCapacity:    41,
		MinimumCapacity:    41,
		MaximumCapacity:    43,
		WarehouseID:        43,
		ProductTypeID:      43,
	}
	repo.On("Update", mock.Anything, objetoAc).Return(nil)
	s := NewService(repo)
	objetoRecuperado, err := s.Update(context.Background(), objetoAc)
	assert.NoError(t, err)
	assert.Equal(t, objetoAc, objetoRecuperado, "Los sections no coinciden")
}

func TestUpdateFail(t *testing.T) {
	repo := new(repoM)
	repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("No se puede modificar el section"))
	s := NewService(repo)
	objetoAc := domain.Section{
		SectionNumber:      41,
		CurrentTemperature: 3,
		MinimumTemperature: 4,
		CurrentCapacity:    41,
		MinimumCapacity:    41,
		MaximumCapacity:    43,
		WarehouseID:        43,
		ProductTypeID:      43,
	}
	objetoRecuperado, err := s.Update(context.Background(), objetoAc)
	assert.Error(t, err)
	assert.Equal(t, domain.Section{}, objetoRecuperado, "Los section no coinciden")
}

func TestDeleteOk(t *testing.T) {
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	s := NewService(repo)
	err := s.Delete(context.Background(), 1)
	assert.NoError(t, err)
}

func TestDeleteFail(t *testing.T) {
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(ErrNotFound)
	s := NewService(repo)
	err := s.Delete(context.Background(), 1)
	assert.Error(t, err)
}
