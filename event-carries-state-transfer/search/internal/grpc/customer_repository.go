package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/search/internal/domain"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository {
	return CustomerRepository{
		client: customerspb.NewCustomersServiceClient(conn),
	}
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	resp, err := r.client.GetCustomer(ctx, &customerspb.GetCustomerRequest{Id: customerID})
	if err != nil {
		return nil, err
	}

	return r.customerToDomain(resp.GetCustomer()), nil
}

func (r CustomerRepository) customerToDomain(customer *customerspb.Customer) *domain.Customer {
	return &domain.Customer{
		ID:   customer.GetId(),
		Name: customer.GetName(),
	}
}
