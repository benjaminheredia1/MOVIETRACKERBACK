package domain

import "time"

type ITEMREOPOSITORY interface {
	Add(item ITEM) error
	GetAll(filters Filters) ([]ITEM, error)
	GetByID(id int) (*ITEM, error)
	MarkAsWatched(id int, rating int, watchedAt time.Time, commentary string) error
	Delete(id int) error
}

type LISTAREOPOSITORY interface {
	Add(list LISTA) (int, error)
	GetAll(filters Filters) ([]LISTA, error)
	GetByID(id int) (*LISTA, error)
	Delete(id int) error
	Update(list LISTA) error
}

type LISTA_ITEMREOPOSITORY interface {
	Add(listItem LISTA_ITEM) error
	GetByListID(listID int) ([]ITEM, error)
	Delete(listID int, itemID int) error
}

type ChatRepository interface {
	GenerateMessage(msg ChatMessage) (string, error)
	GetHistory(sessionID string) ([]ChatMessage, error)
}

type MediaRepository interface {
	Search(query string) ([]MediaResult, error)
	Recomendations() ([]MediaResult, error)
}
