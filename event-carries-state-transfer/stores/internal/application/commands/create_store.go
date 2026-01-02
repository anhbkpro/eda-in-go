package commands

import (
	"context"
	"eda-in-golang/stores/internal/domain"

	"github.com/stackus/errors"
)

type (
	// Command Object is in the same file with Handler
	CreateStore struct {
		Id       string
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
	store, err := domain.CreateStore(cmd.Id, cmd.Name, cmd.Location)
	if err != nil {
		return errors.Wrap(err, "error creating store")
	}

	return errors.Wrap(h.stores.Save(ctx, store), "error saving store")
}
