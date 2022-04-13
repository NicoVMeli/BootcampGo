package section

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound      = errors.New("la seccion no fue encontrada")
	ErrAlreadyExists = errors.New("la seccion ya existe")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Save(ctx context.Context, sec domain.Section) (domain.Section, error)
	Update(ctx context.Context, sec domain.Section) (domain.Section, error)
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, sectionNumber int) bool
}

type service struct {
	sectionRepository Repository
}

func NewService(sectionRepo Repository) Service {
	return &service{
		sectionRepository: sectionRepo,
	}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Section, error) {
	sections, err := s.sectionRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Section, error) {
	sect, err := s.sectionRepository.Get(ctx, id)
	if err != nil {
		return domain.Section{}, err
	}

	return sect, nil
}

func (s *service) Save(ctx context.Context, sect domain.Section) (domain.Section, error) {
	exist := s.sectionRepository.Exists(ctx, sect.SectionNumber)
	if exist {
		return domain.Section{}, ErrAlreadyExists
	}
	newSectionId, err := s.sectionRepository.Save(ctx, sect)
	if err != nil {
		return domain.Section{}, err
	}
	sect.ID = newSectionId
	return sect, nil
}

func (s *service) Update(ctx context.Context, sect domain.Section) (domain.Section, error) {
	err := s.sectionRepository.Update(ctx, sect)
	if err != nil {
		return domain.Section{}, ErrNotFound
	}
	return sect, nil

}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.sectionRepository.Delete(ctx, id)
}

func (s *service) Exists(ctx context.Context, sec int) bool {
	return s.sectionRepository.Exists(ctx, sec)
}
