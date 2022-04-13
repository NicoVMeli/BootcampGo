package seller

import(
	"errors"
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)


type repoM struct{
	mock.Mock
}

func (r *repoM) GetAll(ctx context.Context) ([]domain.Seller, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (r *repoM) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (r *repoM) Exists(ctx context.Context, cid int) bool {
	args := r.Called(ctx, cid)
	return args.Bool(0)
}

func (r *repoM) Save(ctx context.Context, se domain.Seller) (int, error) {
	args := r.Called(ctx, se)
	return args.Int(0), args.Error(1)
}

func (r *repoM) Update(ctx context.Context, se domain.Seller) error {
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

	
	objetoAPersistir := domain.Seller{
		CID:         123,
		CompanyName: "Meli",
		Address: "Bulnes 10",
		Telephone: "123456",
	}

	
	ns, err := s.Save(context.Background(), objetoAPersistir)

	
	assert.NoError(t, err)
	assert.Equal(t, 1, ns.ID)
	assert.Equal(t, 123, ns.CID)
}

func TestCreateConflict(t *testing.T) {
	expectedError := errors.New("el seller ya existe")
	expectedResult := domain.Seller{}
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
	repo.On("GetAll", mock.Anything).Return([]domain.Seller{
		{
			ID:          1,
			CID:         123,
			CompanyName: "Meli",
			Address: "Bulnes 10",
			Telephone: "123456",

		},
		{
			ID:          2,
			CID:         456,
			CompanyName: "Baires Dev",
			Address: " Av Belgrano 1200",
			Telephone: "4833269",
		},
	}, nil)

	s := NewService(repo)

	objetosRecuperados, err := s.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, objetosRecuperados, 2)
}

func TestFindOne( t *testing.T){
	
	expectedResult := domain.Seller{
		 	ID:          1,
	 		CID:         123,
	 	CompanyName: "Meli",
	 	Address: "Bulnes 10",
		Telephone: "123456",
	}

	repo := new(repoM)
	repo.On("Get", mock.Anything, 1).Return(expectedResult, nil)
	s:= NewService(repo)
	ctx := context.Background()
	resul, _ := s.Get(ctx, 1)
	assert.Equal(t, expectedResult, resul)

}


func TestFindOneNoExist(t *testing.T){
	repo := new(repoM)
	repo.On("Get", mock.Anything, mock.Anything).Return(domain.Seller{}, errors.New("No existe un seller con ese id"))

	s := NewService(repo)
	
	objetoRecuperado, err := s.Get(context.Background(), 3)

	obj := domain.Seller{}

	assert.Error(t, err)
	assert.Equal(t,obj,objetoRecuperado, "No existe un seller con ese id")
}


func TestUpdateOk(t *testing.T){
	repo := new(repoM)
	objetoAc := domain.Seller{
		ID:          1,
		CID:         123,
		CompanyName: "Meli",
		Address: "Bulnes 10",
		Telephone: "123456",
	}
	repo.On("Update", mock.Anything, objetoAc).Return(nil)
	s := NewService(repo)

	
	objetoRecuperado, err := s.Update(context.Background(), objetoAc)

	assert.NoError(t, err)
	assert.Equal(t,objetoAc, objetoRecuperado, "Los sellers no coinciden")

	

}

func TestUpdateFail(t *testing.T){
	repo := new(repoM)

	repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("No se puede modificar el seller"))
	s := NewService(repo)

	objetoAc := domain.Seller{
		ID:          1,
		CID:         123,
		CompanyName: "Meli",
		Address: "Bulnes 10",
		Telephone: "123456",
	}
	objetoRecuperado, err := s.Update(context.Background(), objetoAc)

	assert.Error(t, err)
	assert.Equal(t,domain.Seller{}, objetoRecuperado, "Los sellers no coinciden" )
}


func TestDeleteOk(t *testing.T){
	
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(nil)

	s := NewService(repo)
	err := s.Delete(context.Background(), 1)
	assert.NoError(t, err)

}



func TestDeleteFail(t *testing.T){
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(ErrNotFound)

	s := NewService(repo)
	err := s.Delete(context.Background(), 1)
	assert.Error(t, err)

}
