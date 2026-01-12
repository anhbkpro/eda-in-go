package es

import (
	"fmt"

	"eda-in-golang/internal/ddd"
)

type EventApplier interface {
	ApplyEvent(ddd.Event) error
}

type EventCommitter interface {
	CommitEvents()
}

func LoadEvent(v interface{}, event ddd.AggregateEvent) error {
	type loader interface {
		EventApplier
		VersionSetter
	}

	agg, ok := v.(loader)
	if !ok {
		return fmt.Errorf("type %T does not implement loader", v)
	}

	err := agg.ApplyEvent(event)
	if err != nil {
		return err
	}

	agg.setVersion(event.AggregateVersion())

	return nil
}
