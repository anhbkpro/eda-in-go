package domain

type StoreV1 struct {
	Name          string
	Location      string
	Participating bool
}

// SnapshotName implements es.Snapshot
func (StoreV1) SnapshotName() string { return "stores.StoreV1" }
