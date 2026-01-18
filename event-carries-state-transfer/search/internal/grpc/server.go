package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/domain"
	"eda-in-golang/search/searchpb"
)

// An implementation of the SearchServiceServer interface.
type server struct {
	app                                       application.Application
	searchpb.UnimplementedSearchServiceServer // From the doc of type SearchServiceServer interface, embed UnimplementedSearchServiceServer for forward compatibility. It will cause panic if an unimplemented method is ever invoked.
}

func RegisterServer(_ context.Context, app application.Application, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})
	return nil
}

func (s server) SearchOrders(ctx context.Context, req *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error) {
	filters := domain.SearchFilters{
		Next:  req.GetNext(),
		Limit: int(req.GetLimit()),
		Filters: domain.Filters{
			CustomerID: req.GetFilters().GetCustomerId(),
			StoreIDs:   req.GetFilters().GetStoreIds(),
			ProductIDs: req.GetFilters().GetProductIds(),
			MinTotal:   req.GetFilters().GetMinTotal(),
			MaxTotal:   req.GetFilters().GetMaxTotal(),
			Status:     req.GetFilters().GetStatus(),
		},
	}

	if req.GetFilters().GetAfter() != nil {
		filters.Filters.After = req.GetFilters().GetAfter().AsTime()
	}
	if req.GetFilters().GetBefore() != nil {
		filters.Filters.Before = req.GetFilters().GetBefore().AsTime()
	}

	orders, err := s.app.SearchOrders(ctx, filters)
	if err != nil {
		return nil, err
	}

	var next string
	if len(orders) > 0 && len(orders) == filters.Limit {
		// Generate next cursor based on last order
		lastOrder := orders[len(orders)-1]
		cursor := domain.Cursor{
			CreatedAt: lastOrder.CreatedAt,
			ID:        0, // We don't have an ID field, using 0
		}
		if next, err = domain.EncodeCursor(cursor); err != nil {
			return nil, err
		}
	}

	return &searchpb.SearchOrdersResponse{
		Orders: s.ordersToProto(orders),
		Next:   next,
	}, nil
}

func (s server) GetOrder(ctx context.Context, req *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error) {
	order, err := s.app.GetOrder(ctx, application.GetOrder{OrderID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &searchpb.GetOrderResponse{
		Order: s.orderToProto(order),
	}, nil
}

func (s server) ordersToProto(orders []*domain.Order) []*searchpb.Order {
	protoOrders := make([]*searchpb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = s.orderToProto(order)
	}
	return protoOrders
}

func (s server) orderToProto(order *domain.Order) *searchpb.Order {
	items := make([]*searchpb.Order_Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = &searchpb.Order_Item{
			ProductId:   item.ProductID,
			StoreId:     item.StoreID,
			ProductName: item.ProductName,
			StoreName:   item.StoreName,
			Price:       item.Price,
			Quantity:    int64(item.Quantity),
		}
	}

	return &searchpb.Order{
		OrderId:      order.ID,
		CustomerId:   order.CustomerID,
		CustomerName: order.CustomerName,
		Items:        items,
		Total:        order.Total,
		Status:       order.Status,
	}
}
