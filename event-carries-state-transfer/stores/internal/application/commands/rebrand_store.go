package commands

import (
	"context"
	"eda-in-golang/stores/internal/domain"

	"github.com/stackus/errors"
)

// Create Command and Handler for rebrand store
type (
	RebrandStore struct {
		ID   string
		Name string
	}

	RebrandStoreHandler struct {
		stores domain.StoreRepository
	}
)

func NewRebrandStoreHandler(stores domain.StoreRepository) RebrandStoreHandler {
	return RebrandStoreHandler{
		stores: stores,
	}
}

// Implement Handle method
func (h RebrandStoreHandler) RebrandStore(ctx context.Context, cmd RebrandStore) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading store")
	}

	err = store.RebrandStore(cmd.Name)
	if err != nil {
		return errors.Wrap(err, "error rebranding store")
	}

	return errors.Wrap(h.stores.Save(ctx, store), "error saving store")
}
