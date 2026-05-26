package postgres

import (
	"MovieTrackerBack/internal/domain"
	"database/sql"
	"fmt"
	"time"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetAll(filters domain.Filters) ([]domain.ITEM, error) {
	query := `
		SELECT 
			id, tmdb_id, adult, backdrop_path, name, original_name, overview, 
			poster_path, media_type, original_language, popularity, first_air_date, 
			softcore, genre_ids, origin_country, vote_average, vote_count, 
			list_id, status, comentary_user, calification_user, watched_at, added_at 
		FROM items
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando la data: %w", err)
	}
	defer rows.Close()

	var items []domain.ITEM

	for rows.Next() {
		var item domain.ITEM
		err := rows.Scan(
			&item.ID,
			&item.TmdbID,
			&item.Adult,
			&item.BackdropPath,
			&item.Name,
			&item.OriginalName,
			&item.Overview,
			&item.PosterPath,
			&item.MediaType,
			&item.OriginalLanguage,
			&item.Popularity,
			&item.FirstAirDate,
			&item.Softcore,
			&item.GenreIDs,
			&item.OriginCountry,
			&item.VoteAverage,
			&item.VoteCount,
			&item.ListID,
			&item.Status,
			&item.ComentaryUser,
			&item.CalificationUser,
			&item.WatchedAt,
			&item.AddedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error leyendo fila: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ItemRepository) Add(item domain.ITEM) error {

	query := `	INSERT INTO ITEMS (
			tmdb_id, adult, backdrop_path, name, original_name, overview, 
			poster_path, media_type, original_language, popularity, first_air_date, 
			softcore, genre_ids, origin_country, vote_average, vote_count, 
			list_id, status, comentary_user, calification_user, watched_at, added_at
		) 
		VALUES (
			$1, $2, $3, $4, $5, $6, 
			$7, $8, $9, $10, $11, 
			$12, $13, $14, $15, $16, 
			$17, $18, $19, $20, $21, $22
		) RETURNING id`

	var newID int
	err := r.db.QueryRow(query,
		item.TmdbID,
		item.Adult,
		item.BackdropPath,
		item.Name,
		item.OriginalName,
		item.Overview,
		item.PosterPath,
		item.MediaType,
		item.OriginalLanguage,
		item.Popularity,
		item.FirstAirDate,
		item.Softcore,
		item.GenreIDs,
		item.OriginCountry,
		item.VoteAverage,
		item.VoteCount,
		item.ListID,
		item.Status,
		item.ComentaryUser,
		item.CalificationUser,
		item.WatchedAt,
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) GetByID(id int) (*domain.ITEM, error) {
	query := `SELECT 
			id, tmdb_id, adult, backdrop_path, name, original_name, overview, 
			poster_path, media_type, original_language, popularity, first_air_date, 
			softcore, genre_ids, origin_country, vote_average, vote_count, 
			list_id, status, comentary_user, calification_user, watched_at, added_at 
		FROM items WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var item domain.ITEM
	row.Scan(
		&item.ID,
		&item.TmdbID,
		&item.Adult,
		&item.BackdropPath,
		&item.Name,
		&item.OriginalName,
		&item.Overview,
		&item.PosterPath,
		&item.MediaType,
		&item.OriginalLanguage,
		&item.Popularity,
		&item.FirstAirDate,
		&item.Softcore,
		&item.GenreIDs,
		&item.OriginCountry,
		&item.VoteAverage,
		&item.VoteCount,
		&item.ListID,
		&item.Status,
		&item.ComentaryUser,
		&item.CalificationUser,
		&item.WatchedAt,
		&item.AddedAt,
	)
	return &item, nil
}

func (r *ItemRepository) MarkAsWatched(id int, rating int, watchedAt time.Time, commentary string) error {
	query := `UPDATE items SET status = 'watched', calification_user = $1, watched_at = $2, comentary_user = $3 WHERE id = $4`
	_, err := r.db.Exec(query, rating, watchedAt, commentary, id)
	if err != nil {
		return fmt.Errorf("error actualizando el item: %w", err)
	}
	return nil
}

func (r *ItemRepository) Delete(id int) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
