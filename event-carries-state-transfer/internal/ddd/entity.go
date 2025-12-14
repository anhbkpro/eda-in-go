package ddd

type IDer interface {
	ID() string
}

type EntityNamer interface {
	EntityName() string
}

type Entity struct {
	id   string
	name string
}

var _ interface {
	IDer
	EntityNamer
	IDSetter
	NameSetter
} = (*Entity)(nil)

func NewEntity(id string, name string) Entity {
	return Entity{
		id:   id,
		name: name,
	}
}

func (e Entity) ID() string               { return e.id }
func (e Entity) EntityName() string       { return e.name }
func (e Entity) Equals(other Entity) bool { return e.id == other.ID() }

func (e *Entity) setID(id string)     { e.id = id }
func (e *Entity) setName(name string) { e.name = name }
