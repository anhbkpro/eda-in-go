package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/es"
)

const CustomerAggregate = "customers.CustomerAggregate"

type Customer struct {
	es.Aggregate
	Name      string
	SmsNumber string
	Enabled   bool
}

var (
	ErrNameCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the customer name cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer ID cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the customer SMS number cannot be blank")
	ErrCustomerAlreadyEnabled  = errors.Wrap(errors.ErrBadRequest, "the customer is already enabled")
	ErrCustomerAlreadyDisabled = errors.Wrap(errors.ErrBadRequest, "the customer is already disabled")
	ErrCustomerNotAuthorized   = errors.Wrap(errors.ErrForbidden, "the customer is not authorized to perform this action")
)

func NewCustomer(id string) *Customer {
	return &Customer{
		Aggregate: es.NewAggregate(id, CustomerAggregate),
	}
}

func RegisterCustomer(id, name, smsNumber string) (*Customer, error) {
	if name == "" {
		return nil, ErrNameCannotBeBlank
	}
	if smsNumber == "" {
		return nil, ErrSmsNumberCannotBeBlank
	}
	customer := NewCustomer(id)
	customer.Name = name
	customer.SmsNumber = smsNumber
	customer.Enabled = true

	customer.AddEvent(CustomerRegisteredEvent, &CustomerRegistered{
		Customer: customer,
	})

	return customer, nil
}

func (Customer) Key() string {
	return CustomerAggregate
}

func (c *Customer) Authorize( /* TODO authorize what? */ ) error {
	if !c.Enabled {
		return ErrCustomerNotAuthorized
	}

	c.AddEvent(CustomerAuthorizedEvent, &CustomerAuthorized{
		Customer: c,
	})

	return nil
}

func (c *Customer) Enable() error {
	if c.Enabled {
		return ErrCustomerAlreadyEnabled
	}

	c.AddEvent(CustomerEnabledEvent, &CustomerEnabled{
		Customer: c,
	})

	return nil
}

func (c *Customer) Disable() error {
	if !c.Enabled {
		return ErrCustomerAlreadyDisabled
	}

	c.AddEvent(CustomerDisabledEvent, &CustomerDisabled{
		Customer: c,
	})

	return nil
}
