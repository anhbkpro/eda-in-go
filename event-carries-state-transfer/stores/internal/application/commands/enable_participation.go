package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

// Create Command and Handler for enable participation
type (
	EnableParticipation struct {
		ID string
	}

	EnableParticipationHandler struct {
		stores domain.StoreRepository
	}
)

func NewEnableParticipationHandler(stores domain.StoreRepository) EnableParticipationHandler {
	return EnableParticipationHandler{
		stores: stores,
	}
}

// Implement Handle method
func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipation) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading store")
	}

	err = store.EnableParticipation()
	if err != nil {
		return errors.Wrap(err, "error enabling participation")
	}

	return errors.Wrap(h.stores.Save(ctx, store), "error saving store")
}
