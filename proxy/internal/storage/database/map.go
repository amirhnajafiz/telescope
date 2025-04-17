package database

// mapping is a simple in-memory key-value store
type mapping struct {
	entities map[string]string
}

// Add creates a new mapping instance
func (m *mapping) Add(key string, value string) error {
	m.entities[key] = value
	return nil
}

// Get retrieves a value by key
func (m *mapping) Get(key string) (string, error) {
	if value, ok := m.entities[key]; ok {
		return value, nil
	}

	return "", nil
}

// List retrieves all key-value pairs
func (m *mapping) List() ([]string, error) {
	var items []string

	for key := range m.entities {
		items = append(items, key)
	}

	return items, nil
}
