package grpc

import (
	"context"
	"eda-in-golang/depot/depotpb"
	"eda-in-golang/depot/internal/application"
	"eda-in-golang/depot/internal/application/commands"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*server)(nil)

func Register(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(registrar, &server{app: app})
	return nil
}

func (s *server) CreateShoppingList(ctx context.Context, req *depotpb.CreateShoppingListRequest) (*depotpb.CreateShoppingListResponse, error) {
	id := uuid.NewString()

	items := make([]commands.OrderItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		items = append(items, commands.OrderItem{
			StoreID:   item.GetStoreId(),
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}

	err := s.app.CreateShoppingList(ctx, commands.CreateShoppingList{
		ID:      id,
		OrderID: req.GetOrderId(),
		Items:   items,
	})

	return &depotpb.CreateShoppingListResponse{
		Id: id,
	}, err
}

func (s *server) CancelShoppingList(ctx context.Context, req *depotpb.CancelShoppingListRequest) (*depotpb.CancelShoppingListResponse, error) {
	err := s.app.CancelShoppingList(ctx, commands.CancelShoppingList{
		ID: req.GetId(),
	})
	return &depotpb.CancelShoppingListResponse{}, err
}

func (s *server) AssignShoppingList(ctx context.Context, req *depotpb.AssignShoppingListRequest) (*depotpb.AssignShoppingListResponse, error) {
	err := s.app.AssignShoppingList(ctx, commands.AssignShoppingList{
		ID:    req.GetId(),
		BotID: req.GetBotId(),
	})
	return &depotpb.AssignShoppingListResponse{}, err
}

func (s *server) CompleteShoppingList(ctx context.Context, req *depotpb.CompleteShoppingListRequest) (*depotpb.CompleteShoppingListResponse, error) {
	err := s.app.CompleteShoppingList(ctx, commands.CompleteShoppingList{
		ID: req.GetId(),
	})
	return &depotpb.CompleteShoppingListResponse{}, err
}
