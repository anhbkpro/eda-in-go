package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"eda-in-golang/stores/internal/domain"
)

type CatalogRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db *sql.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error {
	// build query
	const query = "INSERT INTO %s (id, store_id, name, description, sku, price) VALUES ($1, $2, $3, $4, $5, $6)"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name, description, sku, price)

	return err
}

func (r CatalogRepository) Rebrand(ctx context.Context, productID, name, description string) error {
	// build query
	const query = "UPDATE %s SET name = $1, description = $2 WHERE id = $3"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), name, description, productID)

	return err
}

func (r CatalogRepository) UpdatePrice(ctx context.Context, productID string, delta float64) error {
	fmt.Printf("[Step 21] CatalogRepo: UPDATE %s SET price = price + %.2f WHERE id = %s\n", r.tableName, delta, productID)

	// build query
	const query = "UPDATE %s SET price = price + $1 WHERE id = $2"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), delta, productID)

	return err
}

func (r CatalogRepository) RemoveProduct(ctx context.Context, productID string) error {
	// build query
	const query = "DELETE FROM %s WHERE id = $1"

	// execute query
	_, err := r.db.ExecContext(ctx, r.table(query), productID)

	return err
}

func (r CatalogRepository) Find(ctx context.Context, productID string) (*domain.CatalogProduct, error) {
	// build query
	const query = "SELECT id, store_id, name, description, sku, price FROM %s WHERE id = $1"

	// create output
	product := &domain.CatalogProduct{}

	// execute query
	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(&product.ID, &product.StoreID, &product.Name, &product.Description, &product.SKU, &product.Price)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r CatalogRepository) GetCatalog(ctx context.Context, storeID string) (products []*domain.CatalogProduct, err error) {
	// build query
	const query = "SELECT id, store_id, name, description, sku, price FROM %s WHERE store_id = $1"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), storeID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing catalog rows")
		}
	}(rows)

	for rows.Next() {
		var product domain.CatalogProduct
		err = rows.Scan(&product.ID, &product.StoreID, &product.Name, &product.Description, &product.SKU, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
