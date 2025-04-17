package database

import "fmt"

// DB is an interface that defines the methods for a database
type DB interface {
	Add(key string, value string) error
	Get(key string) (string, error)
	List() ([]string, error)
}

// NewDB creates a new database instance
func NewDB(dbType string) (DB, error) {
	switch dbType {
	case "map":
		return &mapping{
			Entities: make(map[string]string),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
