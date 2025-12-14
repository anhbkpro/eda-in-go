package grpc

import (
	"context"
	"eda-in-golang/ordering/internal/domain"
	"eda-in-golang/payments/paymentspb"

	"google.golang.org/grpc"
)

type PaymentRepository struct {
	client paymentspb.PaymentsServiceClient
}

var _ domain.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(conn *grpc.ClientConn) domain.PaymentRepository {
	return &PaymentRepository{
		client: paymentspb.NewPaymentsServiceClient(conn),
	}
}

func (r *PaymentRepository) Confirm(ctx context.Context, paymentID string) error {
	_, err := r.client.ConfirmPayment(ctx, &paymentspb.ConfirmPaymentRequest{
		Id: paymentID,
	})
	return err
}
