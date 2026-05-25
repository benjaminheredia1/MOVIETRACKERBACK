package domain

import "time"

type ITEMREOPOSITORY interface {
	Add(item ITEM) error
	GetAll(filters ITEM) ([]ITEM, error)
	GetByID(id int) (*ITEM, error)
	MarkAsWatched(id int, rating int, watchedAt time.Time, commentary string) error
	Delete(id int) error
}

type LISTAREOPOSITORY interface {
	Add(list LISTA) (int, error)
	GetAll(Filters LISTA) ([]LISTA, error)
	GetByID(id int) (*LISTA, error)
	Delete(id int) error
}

type LISTA_ITEMREOPOSITORY interface {
	Add(listItem LISTA_ITEM) error
	GetByListID(listID int) ([]ITEM, error)
	Delete(listID int, itemID int) error
}

// type ChatRepository interface {
// 	SaveMessage(msg ChatMessage) error
// 	GetHistory(sessionID string) ([]ChatMessage, error)
// }

// type CacheRepository interface {
// 	Get(key string) ([]byte, error)
// 	Set(key string, value []byte, ttl time.Duration) error
// 	Delete(key string) error
// }

// type MediaRepository interface {
// 	Search(query string, mediaType string) ([]MediaResult, error)
// 	GetDetail(tmdbID int, mediaType string) (*MediaDetail, error)
// } 
