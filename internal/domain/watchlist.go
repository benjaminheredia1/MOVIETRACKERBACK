package domain

import "time"

// Modelo de lo que se va a guardar en la tabla, con los campos necesarios para mostrar la información relevante al usuario y para marcar como visto o eliminar de la lista.

type ITEM struct {
	ID               int        `json:"id"`
	TmdbID           int        `json:"tmdb_id"`
	Adult            bool       `json:"adult"`
	BackdropPath     string     `json:"backdrop_path"`
	Name             string     `json:"name"`
	OriginalName     string     `json:"original_name"`
	Overview         string     `json:"overview"`
	PosterPath       string     `json:"poster_path"`
	MediaType        string     `json:"media_type"`
	OriginalLanguage string     `json:"original_language"`
	Popularity       float64    `json:"popularity"`
	FirstAirDate     string     `json:"first_air_date"`
	Softcore         bool       `json:"softcore"`
	GenreIDs         string     `json:"genre_ids"`
	OriginCountry    string     `json:"origin_country"`
	VoteAverage      float64    `json:"vote_average"`
	VoteCount        int        `json:"vote_count"`
	ListID           *int       `json:"list_id"`
	Status           string     `json:"status"`
	ComentaryUser    string     `json:"comentary_user"`
	CalificationUser *float64   `json:"calification_user"`
	WatchedAt        *time.Time `json:"watched_at"`
	AddedAt          time.Time  `json:"added_at"`
}

type LISTA struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type LISTA_ITEM struct {
	ID      int       `json:"id"`
	ListID  int       `json:"list_id"`
	ItemID  int       `json:"item_id"`
	AddedAt time.Time `json:"added_at"`
}

type Filters struct {
	Status    string
	MediaType string
	OrderBy   string
	OrderDir  string
}
