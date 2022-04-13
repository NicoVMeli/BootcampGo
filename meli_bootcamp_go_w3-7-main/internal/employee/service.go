package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound           = errors.New("employee not found")
	ErrorEmployeesNoExist = errors.New("The employee doesn't exist")
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Save(ctx context.Context, e domain.Employee) (domain.Employee, error)
	Update(ctx context.Context, id int, e domain.Employee) (domain.Employee, error)
	Delete(ctx context.Context, id int) error
	GetAllReportInboundOrders(ctx context.Context) ([]domain.ReportInboundOrders, error)
	GetReportByEmployeeIdInboundOrders(ctx context.Context, employeeId int) (domain.ReportInboundOrders, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Get(ctx context.Context, id int) (domain.Employee, error) {
	resul, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("No existe el empleado con el id %d", id)
	}
	return resul, err
}

func (s *service) GetAll(ctx context.Context) ([]domain.Employee, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Save(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	exists := s.repository.Exists(ctx, e.CardNumberID)
	if exists {
		return domain.Employee{}, errors.New("El usuario ya existe")
	}
	id, err := s.repository.Save(ctx, e)
	if err != nil {
		return domain.Employee{}, err
	}
	e.ID = id
	return e, nil
}

func (s *service) Update(ctx context.Context, id int, e domain.Employee) (domain.Employee, error) {
	e.ID = id
	resul, _ := s.repository.Get(ctx, id)
	if resul.ID == 0 {
		return domain.Employee{}, fmt.Errorf("No existe el empleado con el id %d", id)
	}
	updateField := updateField(resul, e)
	err := s.repository.Update(ctx, updateField)
	if err != nil {
		return domain.Employee{}, err
	}
	return e, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	resul, _ := s.repository.Get(ctx, id)
	if resul.ID == 0 {
		return fmt.Errorf("No existe el empleado con el id %d", id)
	}
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllReportInboundOrders(ctx context.Context) ([]domain.ReportInboundOrders, error) {
	return s.repository.GetAllReportInboundOrders()
}

func (s *service) GetReportByEmployeeIdInboundOrders(ctx context.Context, employeeId int) (domain.ReportInboundOrders, error) {
	resul, err := s.repository.GetReportByEmployeeIdInboundOrders(employeeId)
	if err != nil {
		return domain.ReportInboundOrders{}, ErrorEmployeesNoExist
	}
	return resul, err
}

func updateField(lastE domain.Employee, newE domain.Employee) domain.Employee {
	if newE.LastName != lastE.LastName && newE.LastName != "" {
		lastE.LastName = newE.LastName
	}
	if newE.FirstName != lastE.FirstName && newE.FirstName != "" {
		lastE.FirstName = newE.FirstName
	}
	if newE.CardNumberID != lastE.CardNumberID && newE.CardNumberID != "" {
		lastE.CardNumberID = newE.CardNumberID
	}
	if newE.WarehouseID != lastE.WarehouseID && newE.WarehouseID != 0 {
		lastE.WarehouseID = newE.WarehouseID
	}
	return lastE
}
