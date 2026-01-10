package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/searchpb"
)

// An implementation of the SearchServiceServer interface.
type server struct {
	app                                       application.Application
	searchpb.UnimplementedSearchServiceServer // From the doc of type SearchServiceServer interface, embed UnimplementedSearchServiceServer for forward compatibility. It will cause panic if an unimplemented method is ever invoked.
}

func RegisterServer(app application.Application, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})
	return nil
}

func (s server) SearchOrders(context.Context, *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error) {
	// TODO implement me
	panic("not implemented")
}

func (s server) GetOrder(context.Context, *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error) {
	// TODO implement me
	panic("not implemented")
}
