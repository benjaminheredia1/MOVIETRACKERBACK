package postgres

import (
	"MovieTrackerBack/internal/domain"
	"database/sql"
)

type ListRepository struct {
	db *sql.DB
}

func NewListaRepository(db *sql.DB) *ListRepository {
	return &ListRepository{db: db}
}

func (r *ListRepository) GetAll(filters domain.Filters) ([]domain.LISTA, error) {
	query := `SELECT id, name, description, created_at  FROM lista`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.LISTA

	for rows.Next() {
		var item domain.LISTA
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ListRepository) GetByID(id int) (*domain.LISTA, error) {
	query := `SELECT id, name, description, created_at FROM lista WHERE id = $1`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var item domain.LISTA
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, err
}

func (r *ListRepository) Add(list domain.LISTA) (int, error) {
	query := `INSERT INTO lista (name, description) VALUES ($1, $2) RETURNING id`
	var id int
	err := r.db.QueryRow(query, list.Name, list.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ListRepository) Delete(id int) error {
	query := `DELETE FROM lista WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ListRepository) Update(list domain.LISTA) error {
	query := `UPDATE lista SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.Exec(query, list.Name, list.Description, list.ID)
	return err
}
