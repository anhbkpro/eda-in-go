package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"eda-in-golang/stores/internal/domain"
)

type MallRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.MallRepository = (*MallRepository)(nil)

func NewMallRepository(tableName string, db *sql.DB) MallRepository {
	return MallRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MallRepository) AddStore(ctx context.Context, storeID, name, location string) error {
	// build query
	const query = "INSERT INTO %s (id, name, location, participating) VALUES ($1, $2, $3, $4)"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name, location, false)

	return err
}

func (r MallRepository) Find(ctx context.Context, storeID string) (*domain.MallStore, error) {
	// build query
	const query = "SELECT id, name, location, participating FROM %s WHERE id = $1"

	// create output
	store := &domain.MallStore{}

	// execute query
	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r MallRepository) RenameStore(ctx context.Context, storeID, name string) error {
	// build query
	const query = "UPDATE %s SET name = $1 WHERE id = $2"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), name, storeID)

	return err
}

func (r MallRepository) SetStoreParticipation(ctx context.Context, storeID string, participating bool) error {
	// build query
	const query = "UPDATE %s SET participating = $1 WHERE id = $2"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), participating, storeID)

	return err
}

func (r MallRepository) All(ctx context.Context) (stores []*domain.MallStore, err error) {
	// build query
	const query = "SELECT id, name, location, participating FROM %s"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing mall rows")
		}
	}(rows)

	for rows.Next() {
		var store domain.MallStore
		err = rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}

	return stores, nil
}

func (r MallRepository) AllParticipating(ctx context.Context) (stores []*domain.MallStore, err error) {
	// build query
	const query = "SELECT id, name, location, participating FROM %s WHERE participating = true"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing mall rows")
		}
	}(rows)

	for rows.Next() {
		var store domain.MallStore
		err = rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}

	return stores, nil
}

func (r MallRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
