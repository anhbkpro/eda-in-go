package ddd

import (
	"fmt"

	"eda-in-golang/internal/registry"
)

type IDSetter interface {
	setID(id string)
}

func WithID(id string) registry.BuildOption {
	return func(v interface{}) error {
		if ider, ok := v.(IDSetter); ok {
			ider.setID(id)
			return nil
		}
		return fmt.Errorf("type %T does not implement IDSetter", v)
	}
}

type NameSetter interface {
	setName(name string)
}

func WithName(name string) registry.BuildOption {
	return func(v interface{}) error {
		if namer, ok := v.(NameSetter); ok {
			namer.setName(name)
			return nil
		}
		return fmt.Errorf("type %T does not implement NameSetter", v)
	}
}
