package storage

type Storage interface {
	List(prefix string) (ListResult, error)
}

type ListResult struct {
	Results []string
}
