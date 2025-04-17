package database

// DB is an interface that defines the methods for a database
type DB interface {
	Add(key string, value string) error
	Get(key string) (string, error)
	List() ([]string, error)
}

// NewDB creates a new database instance
func NewDB() DB {
	return &mapping{
		entities: make(map[string]string),
	}
}
