package application

import (
	"context"
	"eda-in-golang/stores/internal/application/commands"
	"eda-in-golang/stores/internal/application/queries"
	"eda-in-golang/stores/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateStore(ctx context.Context, cmd commands.CreateStore) error
		EnableParticipation(ctx context.Context, cmd commands.EnableParticipation) error
		DisableParticipation(ctx context.Context, cmd commands.DisableParticipation) error
		RebrandStore(ctx context.Context, cmd commands.RebrandStore) error
		AddProduct(ctx context.Context, cmd commands.AddProduct) error
		RebrandProduct(ctx context.Context, cmd commands.RebrandProduct) error
		IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPrice) error
		DecreaseProductPrice(ctx context.Context, cmd commands.DecreaseProductPrice) error
		RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) error
	}

	Queries interface {
		GetStore(ctx context.Context, query queries.GetStore) (*domain.MallStore, error)
		GetStores(ctx context.Context, query queries.GetStores) ([]*domain.MallStore, error)
		GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStores) ([]*domain.MallStore, error)
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.CatalogProduct, error)
		GetProduct(ctx context.Context, query queries.GetProduct) (*domain.CatalogProduct, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddProductHandler
	}

	appQueries struct {
	}
)

var _ App = (*Application)(nil)

func New(products domain.ProductRepository) *Application {
	return &Application{
		appCommands: {
			commands.AddProductHandler: commands.NewAddProductHandler(products),
		},
	}
}
