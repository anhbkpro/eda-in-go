package grpc

import (
	"context"
	"eda-in-golang/customers/customerspb"
	"eda-in-golang/ordering/internal/domain"

	"google.golang.org/grpc"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) domain.CustomerRepository {
	return &CustomerRepository{
		client: customerspb.NewCustomersServiceClient(conn),
	}
}

func (r *CustomerRepository) Authorize(ctx context.Context, customerID string) error {
	_, err := r.client.AuthorizeCustomer(ctx, &customerspb.AuthorizeCustomerRequest{
		Id: customerID,
	})
	return err
}
