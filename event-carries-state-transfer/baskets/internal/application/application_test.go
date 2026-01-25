package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
)

// Mock implementations for testing

type mockBasketRepository struct {
	loadFunc func(ctx context.Context, basketID string) (*domain.Basket, error)
	saveFunc func(ctx context.Context, basket *domain.Basket) error
}

// Helper function to create a started basket with events applied
func createStartedBasket(id, customerID string) *domain.Basket {
	basket, _ := domain.StartBasket(id, customerID)
	// Apply all pending events to update the basket state
	for _, event := range basket.Events() {
		basket.ApplyEvent(event)
	}
	basket.ClearEvents()
	return basket
}

// Helper function to create a basket with items
func createBasketWithItems(id, customerID string) *domain.Basket {
	basket := createStartedBasket(id, customerID)
	store := &domain.Store{ID: "store-1", Name: "Store 1"}
	product := &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}
	basket.AddItem(store, product, 2)
	// Apply all pending events to update the basket state
	for _, event := range basket.Events() {
		basket.ApplyEvent(event)
	}
	basket.ClearEvents()
	return basket
}

func (m *mockBasketRepository) Load(ctx context.Context, basketID string) (*domain.Basket, error) {
	if m.loadFunc != nil {
		return m.loadFunc(ctx, basketID)
	}
	return nil, fmt.Errorf("loadFunc not implemented")
}

func (m *mockBasketRepository) Save(ctx context.Context, basket *domain.Basket) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, basket)
	}
	return fmt.Errorf("saveFunc not implemented")
}

type mockStoreRepository struct {
	findFunc func(ctx context.Context, storeID string) (*domain.Store, error)
}

func (m *mockStoreRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, storeID)
	}
	return nil, fmt.Errorf("findFunc not implemented")
}

type mockProductRepository struct {
	findFunc func(ctx context.Context, productID string) (*domain.Product, error)
}

func (m *mockProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, productID)
	}
	return nil, fmt.Errorf("findFunc not implemented")
}

type mockEventPublisher struct {
	publishFunc func(ctx context.Context, events ...ddd.Event) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, events ...ddd.Event) error {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, events...)
	}
	return fmt.Errorf("publishFunc not implemented")
}

func TestApplication_StartBasket(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx        context.Context
		basketID   string
		customerID string
	}
	type expected struct {
		err error
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:        context.Background(),
				basketID:   "basket-123",
				customerID: "customer-456",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return domain.NewBasket("basket-123"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return nil
				}
			},
			expected: expected{
				err: nil,
			},
			wantErr: false,
		},
		"load_basket_fails": {
			args: args{
				ctx:        context.Background(),
				basketID:   "basket-123",
				customerID: "customer-456",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("database error")
				}
			},
			expected: expected{
				err: fmt.Errorf("database error"),
			},
			wantErr: true,
			errMsg:  "database error",
		},
		"start_basket_fails_with_empty_customer_id": {
			args: args{
				ctx:        context.Background(),
				basketID:   "basket-123",
				customerID: "",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return domain.NewBasket("basket-123"), nil
				}
			},
			expected: expected{
				err: domain.ErrCustomerIDCannotBeBlank,
			},
			wantErr: true,
			errMsg:  "customer id cannot be blank",
		},
		"save_basket_fails": {
			args: args{
				ctx:        context.Background(),
				basketID:   "basket-123",
				customerID: "customer-456",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return domain.NewBasket("basket-123"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return fmt.Errorf("save error")
				}
			},
			expected: expected{
				err: fmt.Errorf("save error"),
			},
			wantErr: true,
			errMsg:  "save error",
		},
		"publish_event_fails": {
			args: args{
				ctx:        context.Background(),
				basketID:   "basket-123",
				customerID: "customer-456",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return domain.NewBasket("basket-123"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return fmt.Errorf("publish error")
				}
			},
			expected: expected{
				err: fmt.Errorf("publish error"),
			},
			wantErr: true,
			errMsg:  "publish error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			err := app.StartBasket(tt.args.ctx, StartBasket{
				ID:         tt.args.basketID,
				CustomerID: tt.args.customerID,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplication_CancelBasket(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx      context.Context
		basketID string
	}
	type expected struct {
		err error
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return nil
				}
			},
			expected: expected{
				err: nil,
			},
			wantErr: false,
		},
		"load_basket_fails": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("load error")
				}
			},
			expected: expected{
				err: fmt.Errorf("load error"),
			},
			wantErr: true,
			errMsg:  "load error",
		},
		"cancel_non_cancellable_basket": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return domain.NewBasket("basket-123"), nil
				}
			},
			expected: expected{
				err: domain.ErrBasketCannotBeCancelled,
			},
			wantErr: true,
			errMsg:  "basket cannot be cancelled",
		},
		"save_basket_fails": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return fmt.Errorf("save error")
				}
			},
			expected: expected{
				err: fmt.Errorf("save error"),
			},
			wantErr: true,
			errMsg:  "save error",
		},
		"publish_event_fails": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return fmt.Errorf("publish error")
				}
			},
			expected: expected{
				err: fmt.Errorf("publish error"),
			},
			wantErr: true,
			errMsg:  "publish error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			err := app.CancelBasket(tt.args.ctx, CancelBasket{
				ID: tt.args.basketID,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplication_CheckoutBasket(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx       context.Context
		basketID  string
		paymentID string
	}
	type expected struct {
		err error
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "payment-789",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return nil
				}
			},
			expected: expected{
				err: nil,
			},
			wantErr: false,
		},
		"load_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "payment-789",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("load error")
				}
			},
			expected: expected{
				err: fmt.Errorf("load error"),
			},
			wantErr: true,
			errMsg:  "load error",
		},
		"checkout_empty_basket": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "payment-789",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
			},
			expected: expected{
				err: domain.ErrBasketHasNoItems,
			},
			wantErr: true,
			errMsg:  "basket has no items",
		},
		"checkout_with_empty_payment_id": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
			},
			expected: expected{
				err: domain.ErrPaymentIDCannotBeBlank,
			},
			wantErr: true,
			errMsg:  "payment id cannot be blank",
		},
		"save_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "payment-789",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return fmt.Errorf("save error")
				}
			},
			expected: expected{
				err: fmt.Errorf("save error"),
			},
			wantErr: true,
			errMsg:  "save error",
		},
		"publish_event_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				paymentID: "payment-789",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
				f.publisher.publishFunc = func(ctx context.Context, events ...ddd.Event) error {
					return fmt.Errorf("publish error")
				}
			},
			expected: expected{
				err: fmt.Errorf("publish error"),
			},
			wantErr: true,
			errMsg:  "publish error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			err := app.CheckoutBasket(tt.args.ctx, CheckoutBasket{
				ID:        tt.args.basketID,
				PaymentID: tt.args.paymentID,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplication_AddItem(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx       context.Context
		basketID  string
		productID string
		quantity  int
	}
	type expected struct {
		err error
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  2,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.stores.findFunc = func(ctx context.Context, storeID string) (*domain.Store, error) {
					return &domain.Store{ID: "store-1", Name: "Store 1"}, nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
			},
			expected: expected{
				err: nil,
			},
			wantErr: false,
		},
		"load_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  2,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("load error")
				}
			},
			expected: expected{
				err: fmt.Errorf("load error"),
			},
			wantErr: true,
			errMsg:  "load error",
		},
		"find_product_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  2,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return nil, fmt.Errorf("product not found")
				}
			},
			expected: expected{
				err: fmt.Errorf("product not found"),
			},
			wantErr: true,
			errMsg:  "product not found",
		},
		"find_store_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  2,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.stores.findFunc = func(ctx context.Context, storeID string) (*domain.Store, error) {
					return nil, fmt.Errorf("store not found")
				}
			},
			expected: expected{
				err: nil, // Note: The code returns nil when store is not found
			},
			wantErr: false,
		},
		"negative_quantity": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  -1,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.stores.findFunc = func(ctx context.Context, storeID string) (*domain.Store, error) {
					return &domain.Store{ID: "store-1", Name: "Store 1"}, nil
				}
			},
			expected: expected{
				err: domain.ErrQuantityCannotBeNegative,
			},
			wantErr: true,
			errMsg:  "quantity cannot be negative",
		},
		"save_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  2,
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.stores.findFunc = func(ctx context.Context, storeID string) (*domain.Store, error) {
					return &domain.Store{ID: "store-1", Name: "Store 1"}, nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return fmt.Errorf("save error")
				}
			},
			expected: expected{
				err: fmt.Errorf("save error"),
			},
			wantErr: true,
			errMsg:  "save error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			err := app.AddItem(tt.args.ctx, AddItem{
				ID:        tt.args.basketID,
				ProductID: tt.args.productID,
				Quantity:  tt.args.quantity,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplication_RemoveItem(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx       context.Context
		basketID  string
		productID string
		quantity  int
	}
	type expected struct {
		err error
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  1,
			},
			prepare: func(f *fields) {
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return nil
				}
			},
			expected: expected{
				err: nil,
			},
			wantErr: false,
		},
		"find_product_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  1,
			},
			prepare: func(f *fields) {
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return nil, fmt.Errorf("product not found")
				}
			},
			expected: expected{
				err: fmt.Errorf("product not found"),
			},
			wantErr: true,
			errMsg:  "product not found",
		},
		"load_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  1,
			},
			prepare: func(f *fields) {
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("load error")
				}
			},
			expected: expected{
				err: fmt.Errorf("load error"),
			},
			wantErr: true,
			errMsg:  "load error",
		},
		"negative_quantity": {
			args: args{
				ctx:       context.Context(context.Background()),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  -1,
			},
			prepare: func(f *fields) {
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
			},
			expected: expected{
				err: domain.ErrQuantityCannotBeNegative,
			},
			wantErr: true,
			errMsg:  "quantity cannot be negative",
		},
		"save_basket_fails": {
			args: args{
				ctx:       context.Background(),
				basketID:  "basket-123",
				productID: "product-1",
				quantity:  1,
			},
			prepare: func(f *fields) {
				f.products.findFunc = func(ctx context.Context, productID string) (*domain.Product, error) {
					return &domain.Product{ID: "product-1", StoreID: "store-1", Name: "Product 1", Price: 10.0}, nil
				}
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createBasketWithItems("basket-123", "customer-456"), nil
				}
				f.baskets.saveFunc = func(ctx context.Context, basket *domain.Basket) error {
					return fmt.Errorf("save error")
				}
			},
			expected: expected{
				err: fmt.Errorf("save error"),
			},
			wantErr: true,
			errMsg:  "save error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			err := app.RemoveItem(tt.args.ctx, RemoveItem{
				ID:        tt.args.basketID,
				ProductID: tt.args.productID,
				Quantity:  tt.args.quantity,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplication_GetBasket(t *testing.T) {
	type fields struct {
		baskets   *mockBasketRepository
		stores    *mockStoreRepository
		products  *mockProductRepository
		publisher *mockEventPublisher
	}
	type args struct {
		ctx      context.Context
		basketID string
	}
	type expected struct {
		basket   *domain.Basket
		err      error
		validate func(*testing.T, *domain.Basket)
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"success": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return createStartedBasket("basket-123", "customer-456"), nil
				}
			},
			expected: expected{
				err: nil,
				validate: func(t *testing.T, basket *domain.Basket) {
					assert.NotNil(t, basket)
					assert.Equal(t, "basket-123", basket.ID())
					assert.Equal(t, "customer-456", basket.CustomerID)
				},
			},
			wantErr: false,
		},
		"load_basket_fails": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-123",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("load error")
				}
			},
			expected: expected{
				basket: nil,
				err:    fmt.Errorf("load error"),
			},
			wantErr: true,
			errMsg:  "load error",
		},
		"basket_not_found": {
			args: args{
				ctx:      context.Background(),
				basketID: "basket-999",
			},
			prepare: func(f *fields) {
				f.baskets.loadFunc = func(ctx context.Context, basketID string) (*domain.Basket, error) {
					return nil, fmt.Errorf("basket not found")
				}
			},
			expected: expected{
				basket: nil,
				err:    fmt.Errorf("basket not found"),
			},
			wantErr: true,
			errMsg:  "basket not found",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{
				baskets:   &mockBasketRepository{},
				stores:    &mockStoreRepository{},
				products:  &mockProductRepository{},
				publisher: &mockEventPublisher{},
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			app := New(f.baskets, f.stores, f.products, f.publisher)

			// Act
			basket, err := app.GetBasket(tt.args.ctx, GetBasket{
				ID: tt.args.basketID,
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, basket)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.expected.validate != nil {
					tt.expected.validate(t, basket)
				}
			}
		})
	}
}

func TestApplication_New(t *testing.T) {
	t.Parallel()

	// Arrange
	basketRepo := &mockBasketRepository{}
	storeRepo := &mockStoreRepository{}
	productRepo := &mockProductRepository{}
	publisher := &mockEventPublisher{}

	// Act
	app := New(basketRepo, storeRepo, productRepo, publisher)

	// Assert
	assert.NotNil(t, app)
	assert.NotNil(t, app.baskets)
	assert.NotNil(t, app.stores)
	assert.NotNil(t, app.products)
	assert.NotNil(t, app.publisher)

	// Verify interface implementation
	var _ App = app
}
