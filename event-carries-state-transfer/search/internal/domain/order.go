package domain

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Order struct {
	ID           string
	CustomerID   string
	CustomerName string
	Items        []Item
	Total        float64
	Status       string
	CreatedAt    time.Time
}

type Item struct {
	ProductID   string
	StoreID     string
	ProductName string
	StoreName   string
	Price       float64
	Quantity    int
}

type Filters struct {
	CustomerID string    `json:"customer_id"`
	After      time.Time `json:"after"`
	Before     time.Time `json:"before"`
	StoreIDs   []string  `json:"store_ids"`
	ProductIDs []string  `json:"product_ids"`
	MinTotal   float64   `json:"min_total"`
	MaxTotal   float64   `json:"max_total"`
	Status     string    `json:"status"`
}

type SearchFilters struct {
	Filters Filters `json:"filters"`
	Next    string  `json:"next"`
	Limit   int     `json:"limit"`
}

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

func EncodeCursor(c Cursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Cursor{}, err
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return Cursor{}, err
	}
	return c, nil
}
