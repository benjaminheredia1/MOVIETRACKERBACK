package application

import (
	"MovieTrackerBack/internal/domain"
	"errors"
	"time"
)

type ItemsListService struct {
    repo domain.ITEMREOPOSITORY
}

func NewItemsListService(repo domain.ITEMREOPOSITORY) *ItemsListService {
    return &ItemsListService{repo: repo}
}

func (s *ItemsListService) AddItem(item domain.ITEM) error {
    if item.TmdbID == 0 {
        return errors.New("tmdb_id es requerido")
    }
    if item.MediaType != "movie" && item.MediaType != "tv" {
        return errors.New("media_type debe ser movie o tv")
    }
    return s.repo.Add(item)
}

func (s *ItemsListService) GetAll(filters domain.Filters) ([]domain.ITEM, error) {
    return s.repo.GetAll(filters)
}

func (s *ItemsListService) GetByID(id int) (*domain.ITEM, error) {
	
    if id == 0 {
        return nil, errors.New("id es requerido")
    }
    return s.repo.GetByID(id)
}

func (s *ItemsListService) MarkAsWatched(id int, rating int, commentary string) error {
    if rating < 1 || rating > 10 {
        return errors.New("la calificacion debe ser entre 1 y 10")
    }
    return s.repo.MarkAsWatched(id, rating, time.Now(), commentary)
}

func (s *ItemsListService) Delete(id int) error {
    if id == 0 {
        return errors.New("id es requerido")
    }
    return s.repo.Delete(id)
}