// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"github.com/dgraph-io/ristretto"
	ristrettostore "github.com/eko/gocache/store/ristretto/v4"
)

// RistrettoStore is a function which provides a ristrettostore.RistrettoStore instance.
func RistrettoStore() (*ristrettostore.RistrettoStore, error) {
	const (
		// numCounters is the number of counters that the cache will use.
		numCounters = 1e7

		// maxCost is the maximum cost that the cache will use.
		maxCost = 1 << 30

		// bufferItems is the number of items that the cache will buffer.
		bufferItems = 64
	)

	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		return nil, err
	}

	return ristrettostore.NewRistretto(c), nil
}
