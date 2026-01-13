package commands

import (
	"context"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

type (
	// Command Object is in the same file with Handler
	CreateStore struct {
		ID       string
		Name     string
		Location string
	}

	CreateStoreHandler struct {
		stores domain.StoreRepository
	}
)

func NewCreateStoreHandler(stores domain.StoreRepository) CreateStoreHandler {
	return CreateStoreHandler{
		stores: stores,
	}
}

// This is the handler method for the CreateStore command
// Flow: Handler -> Command -> Domain -> Event Bus
func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStore) error {
	fmt.Printf("[Step 1] CreateStoreHandler: processing command for store ID %s\n", cmd.ID)
	fmt.Printf("[Step 2] CreateStoreHandler → Domain.CreateStore: creating store '%s' at '%s'\n", cmd.Name, cmd.Location)

	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location)
	if err != nil {
		fmt.Printf("[Step 2.1] CreateStoreHandler → Domain.CreateStore: ERROR creating store: %v\n", err)
		return errors.Wrap(err, "error creating store")
	}
	fmt.Printf("[Step 2.2] CreateStoreHandler → Domain.CreateStore: store created successfully with %d pending events\n", len(store.Events()))

	fmt.Printf("[Step 3] CreateStoreHandler → StoreRepository.Save: saving store to repository\n")
	err = h.stores.Save(ctx, store)
	if err != nil {
		fmt.Printf("[Step 3.1] CreateStoreHandler → StoreRepository.Save: ERROR saving store: %v\n", err)
		return errors.Wrap(err, "error saving store")
	}
	fmt.Printf("[Step 4] CreateStoreHandler: store saved successfully, command completed\n")

	return nil
}
