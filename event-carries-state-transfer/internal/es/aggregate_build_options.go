package es

import (
	"fmt"

	"eda-in-golang/internal/registry"
)

type VersionSetter interface {
	setVersion(int)
}

func WithVersion(version int) registry.BuildOption {
	return func(v interface{}) error {
		if versioner, ok := v.(VersionSetter); ok {
			versioner.setVersion(version)
			return nil
		}

		return fmt.Errorf("type %T does not implement VersionSetter", v)
	}
}
