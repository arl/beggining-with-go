package shorten

// Store is a simple key value store.
type Store interface {
	// Add saves a new value in the store and returns the unique key.
	Add(value string) uint64

	// Load loads the value associated with a key (or false if it doesn't exist).
	Load(key uint64) (value string, ok bool)
}
