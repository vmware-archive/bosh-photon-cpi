package disk

type Creator interface {
	Create(size int, tenantId string, projectId string, flavorName string) (Disk, error)
}

type Finder interface {
	Find(id string) (Disk, bool, error)
}

type Disk interface {
	ID() string

	Delete() error
}
