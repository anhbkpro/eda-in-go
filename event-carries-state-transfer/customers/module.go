package customers

import (
	"context"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/customers/internal/application"
	"eda-in-golang/customers/internal/grpc"
	"eda-in-golang/customers/internal/handlers"
	"eda-in-golang/customers/internal/logging"
	"eda-in-golang/customers/internal/postgres"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
)

type Module struct{}

func (Module) Name() string {
	return "customers"
}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {
	// setup Driven adapters (Driven: Who your system depends on)
	reg := registry.New()
	if err = customerspb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())
	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)
	domainEventHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		handlers.NewDomainEventHandlers(eventStream),
		"DomainEvents", mono.Logger(),
	)
	commandHandlers := logging.LogCommandHandlerAccess[ddd.Command](
		handlers.NewCommandHandlers(app),
		"Commands", mono.Logger(),
	)

	if err = grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	handlers.RegisterDomainEventHandlers(domainEventHandlers, domainDispatcher)
	err = handlers.RegisterCommandHandlers(commandStream, commandHandlers)
	if err != nil {
		return err
	}

	return nil
}
