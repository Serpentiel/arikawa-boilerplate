// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/container"
	"github.com/eko/gocache/lib/v4/cache"
	ristrettostore "github.com/eko/gocache/store/ristretto/v4"
)

// Cache is a function which provides a container.Cache instance.
func Cache(rs *ristrettostore.RistrettoStore) *container.Cache {
	return container.NewCache(
		cache.New[any](rs),
		cache.New[string](rs),
	)
}
