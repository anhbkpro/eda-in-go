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
		commands.CreateStoreHandler
		commands.DecreaseProductPriceHandler
		commands.IncreaseProductPriceHandler
		commands.EnableParticipationHandler
		commands.DisableParticipationHandler
		commands.RebrandStoreHandler
		commands.RebrandProductHandler
		commands.RemoveProductHandler
	}

	appQueries struct {
		queries.GetCatalogHandler
		queries.GetStoreHandler
		queries.GetParticipatingStoresHandler
		queries.GetProductHandler
		queries.GetStoresHandler
	}
)

var _ App = (*Application)(nil)

func New(
	stores domain.StoreRepository,
	products domain.ProductRepository,
	catalog domain.CatalogRepository,
	mall domain.MallRepository,
) *Application {
	return &Application{
		appCommands: appCommands{
			AddProductHandler:           commands.NewAddProductHandler(products),
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores),
			DecreaseProductPriceHandler: commands.NewDecreaseProductPriceHandler(products),
			IncreaseProductPriceHandler: commands.NewIncreaseProductPriceHandler(products),
			EnableParticipationHandler:  commands.NewEnableParticipationHandler(stores),
			DisableParticipationHandler: commands.NewDisableParticipationHandler(stores),
			RebrandStoreHandler:         commands.NewRebrandStoreHandler(stores),
			RebrandProductHandler:       commands.NewRebrandProductHandler(products),
			RemoveProductHandler:        commands.NewRemoveProductHandler(products),
		},
		appQueries: appQueries{
			GetCatalogHandler:             queries.NewGetCatalogHandler(catalog),
			GetStoreHandler:               queries.NewGetStoreHandler(mall),
			GetParticipatingStoresHandler: queries.NewGetParticipatingStoresHandler(mall),
			GetProductHandler:             queries.NewGetProductHandler(catalog),
			GetStoresHandler:              queries.NewGetStoresHandler(mall),
		},
	}
}
