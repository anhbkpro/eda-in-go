package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"eda-in-golang/search/internal/domain"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackus/errors"
)

type OrderRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db *sql.DB) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Add(ctx context.Context, order *domain.Order) error {
	const query = `INSERT INTO %s (
		order_id, customer_id, customer_name,
		items, status, product_ids, store_ids,
		created_at) VALUES (
		$1, $2, $3,
		$4, $5, $6, $7
		$8
		)`

	items, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	productIDs := make(IDArray, 0, len(order.Items))
	storeMap := make(map[string]struct{})
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
		storeMap[item.StoreID] = struct{}{}
	}
	storeIDs := make(IDArray, 0, len(storeMap))
	for storeID := range storeMap {
		storeIDs = append(storeIDs, storeID)
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		order.ID, order.CustomerID, order.CustomerName,
		items, order.Status, productIDs, storeIDs,
		order.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) Search(ctx context.Context, search *domain.SearchFilters) ([]*domain.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (r OrderRepository) Get(ctx context.Context, orderID string) (*domain.Order, error) {
	const query = `SELECT customer_id, customer_name, items, status, created_at FROM %s WHERE order_id = $1`

	order := &domain.Order{
		ID: orderID,
	}

	var itemData []byte
	err := r.db.QueryRowContext(ctx, r.table(query)).Scan(&order.CustomerID, &order.CustomerName, &itemData, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	var items []domain.Item
	err = json.Unmarshal(itemData, &items)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r OrderRepository) Update(ctx context.Context, orderID string, updater func(*domain.Order) error) error {
	const query = `UPDATE %s SET customer_id = $1, customer_name = $2, items = $3, status = $4, created_at = $5 WHERE order_id = $6`

	order, err := r.Get(ctx, orderID)
	if err != nil {
		return err
	}

	if err := updater(order); err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.CustomerID, order.CustomerName, order.Items, order.Status, order.CreatedAt, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

type IDArray []string

func (a *IDArray) Scan(src any) error {
	var sep = []byte(",")

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.ErrInvalidArgument.Msgf("IDArray: unsupported type: %T", src)
	}

	ids := make([]string, bytes.Count(data, sep))
	for i, id := range bytes.Split(data, sep) {
		ids[i] = string(id)
	}

	*a = ids

	return nil
}

func (a IDArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}

	// unsafe way to do this; assumption is all ids are UUIDs
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}
