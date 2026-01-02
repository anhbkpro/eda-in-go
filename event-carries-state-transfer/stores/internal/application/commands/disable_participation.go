package commands

import (
	"context"
	"eda-in-golang/stores/internal/domain"

	"github.com/stackus/errors"
)

// Create Command and Handler for disable participation
type (
	DisableParticipation struct {
		ID string
	}

	DisableParticipationHandler struct {
		stores domain.StoreRepository
	}
)

func NewDisableParticipationHandler(stores domain.StoreRepository) DisableParticipationHandler {
	return DisableParticipationHandler{
		stores: stores,
	}
}

// Implement Handle method
func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading store")
	}

	err = store.DisableParticipation()
	if err != nil {
		return errors.Wrap(err, "error disabling participation")
	}

	return errors.Wrap(h.stores.Save(ctx, store), "error saving store")
}
