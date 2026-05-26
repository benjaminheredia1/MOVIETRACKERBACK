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

	query := `INSERT INTO items (nombre, descripcion) VALUES ($1, $2)`

	_, err := r.db.Exec(query, item.Name, item.Overview)

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) GetByID(id int) (*domain.ITEM, error) {
	return nil, fmt.Errorf("metodo GetByID no implementado")
}

func (r *ItemRepository) MarkAsWatched(id int, rating int, watchedAt time.Time, commentary string) error {
	return fmt.Errorf("metodo MarkAsWatched no implementado")
}

func (r *ItemRepository) Delete(id int) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
