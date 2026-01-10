package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
)

const ShoppingListAggregate = "depot.ShoppingList"

var (
	ErrShoppingCannotBeCanceled  = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be canceled")
	ErrShoppingCannotBeAssigned  = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be assigned")
	ErrShoppingCannotBeCompleted = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be completed")
)

type ShoppingList struct {
	ddd.Aggregate
	OrderID       string
	Stops         Stops
	AssignedBotID string
	Status        ShoppingListStatus
}

func NewShoppingList(id string) *ShoppingList {
	return &ShoppingList{
		Aggregate: ddd.NewAggregate(id, ShoppingListAggregate),
	}
}

func CreateShoppingList(id, orderID string) *ShoppingList {
	shoppingList := NewShoppingList(id)
	shoppingList.OrderID = orderID
	shoppingList.Status = ShoppingListIsAvailable
	shoppingList.Stops = make(Stops)

	shoppingList.AddEvent(ShoppingListCreatedEvent, &ShoppingListCreated{ShoppingList: shoppingList})

	return shoppingList
}

func (ShoppingList) Key() string { return ShoppingListAggregate }

func (s *ShoppingList) AddItem(store *Store, product *Product, quantity int) error {
	if _, exists := s.Stops[store.ID]; !exists {
		s.Stops[store.ID] = &Stop{
			StoreName:     store.Name,
			StoreLocation: store.Location,
			Items:         make(Items),
		}
	}

	return s.Stops[store.ID].AddItem(product, quantity)
}

func (sl ShoppingList) isCancelable() bool {
	switch sl.Status {
	case ShoppingListIsAvailable, ShoppingListIsAssigned, ShoppingListIsActive:
		return true
	default:
		return false
	}
}

func (s *ShoppingList) Cancel() error {
	if !s.isCancelable() {
		return ErrShoppingCannotBeCanceled
	}

	s.Status = ShoppingListIsCanceled

	s.AddEvent(ShoppingListCanceledEvent, &ShoppingListCanceled{ShoppingList: s})

	return nil
}

func (s ShoppingList) isAssignable() bool {
	return s.Status == ShoppingListIsAvailable
}

func (s *ShoppingList) Assign(botID string) error {
	if !s.isAssignable() {
		return ErrShoppingCannotBeAssigned
	}

	s.AssignedBotID = botID
	s.Status = ShoppingListIsAssigned

	s.AddEvent(ShoppingListAssignedEvent, &ShoppingListAssigned{ShoppingList: s, BotID: botID})

	return nil
}

func (s ShoppingList) isCompletable() bool {
	return s.Status == ShoppingListIsAssigned
}

func (s *ShoppingList) Complete() error {
	if !s.isCompletable() {
		return ErrShoppingCannotBeCompleted
	}

	s.Status = ShoppingListIsCompleted

	s.AddEvent(ShoppingListCompletedEvent, &ShoppingListCompleted{ShoppingList: s})

	return nil
}
