package application

import (
	"MovieTrackerBack/internal/domain"
	"errors"
)

type ListService struct {
	repo domain.LISTAREOPOSITORY
}

func NewListaService(repo domain.LISTAREOPOSITORY) *ListService {
	return &ListService{repo: repo}
}

func (s *ListService) GetAll() ([]domain.LISTA, error) {
	return s.repo.GetAll(domain.Filters{})
}

func (s *ListService) GetByID(id int) (*domain.LISTA, error) {
	if id == 0 {
		return nil, errors.New("id es requerido")
	}
	return s.repo.GetByID(id)
}

func (s *ListService) Add(list domain.LISTA) (int, error) {
	if list.Name == "" {
		return 0, errors.New("name es requerido")
	}
	return s.repo.Add(list)
}

func (s *ListService) Delete(id int) error {
	if id == 0 {
		return errors.New("id es requerido")
	}
	return s.repo.Delete(id)
}

func (s *ListService) Update(list domain.LISTA) error {
	if list.ID == 0 {
		return errors.New("id es requerido")
	}
	if list.Name == "" {
		return errors.New("name es requerido")
	}
	return s.repo.Update(list)
}
