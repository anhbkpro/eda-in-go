// Package domain contains the domain model for the stores module.
//
// For gRPC request examples, see scripts/grpc-requests.sh
package domain

// Product domain event names used for event sourcing.
// These constants are used as event identifiers when storing and retrieving events.
const (
	// ProductAddedEvent is emitted when a new product is added to a store's catalog.
	ProductAddedEvent = "stores.ProductAdded"

	// ProductRebrandedEvent is emitted when a product's name or description is updated.
	ProductRebrandedEvent = "stores.ProductRebranded"

	// ProductPriceIncreasedEvent is emitted when a product's price is increased.
	ProductPriceIncreasedEvent = "stores.ProductPriceIncreased"

	// ProductPriceDecreasedEvent is emitted when a product's price is decreased.
	ProductPriceDecreasedEvent = "stores.ProductPriceDecreased"

	// ProductRemovedEvent is emitted when a product is removed from a store's catalog.
	ProductRemovedEvent = "stores.ProductRemoved"
)

// ProductAdded represents the event payload when a new product is added to a store.
// This event carries the full state of the product at creation time.
type ProductAdded struct {
	StoreID     string  // StoreID is the unique identifier of the store this product belongs to.
	Name        string  // Name is the display name of the product.
	Description string  // Description provides details about the product.
	SKU         string  // SKU is the Stock Keeping Unit identifier.
	Price       float64 // Price is the product's price in the default currency.
}

// Key implements registry.Registerable interface.
// Returns the event name used for serialization/deserialization.
func (ProductAdded) Key() string {
	return ProductAddedEvent
}

// ProductRebranded represents the event payload when a product is rebranded.
// This includes changes to the product's name and/or description.
type ProductRebranded struct {
	Name        string // Name is the new display name of the product.
	Description string // Description is the new product description.
}

// Key implements registry.Registerable interface.
// Returns the event name used for serialization/deserialization.
func (ProductRebranded) Key() string {
	return ProductRebrandedEvent
}

// ProductPriceChanged represents the event payload for price changes.
// This struct is used for both ProductPriceIncreasedEvent and ProductPriceDecreasedEvent.
// The Delta value is positive for increases and negative for decreases.
type ProductPriceChanged struct {
	Delta float64 // Delta is the amount the price changed (positive for increase, negative for decrease).
}

// ProductRemoved represents the event payload when a product is removed from a store.
// This is a marker event with no additional payload data.
type ProductRemoved struct{}

// Key implements registry.Registerable interface.
// Returns the event name used for serialization/deserialization.
func (ProductRemoved) Key() string {
	return ProductRemovedEvent
}
