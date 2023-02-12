// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/cachecontainer"
	"github.com/eko/gocache/lib/v4/cache"
	ristrettostore "github.com/eko/gocache/store/ristretto/v4"
)

// CacheContainer is a function which provides a cachecontainer.CacheContainer instance.
func CacheContainer(rs *ristrettostore.RistrettoStore) *cachecontainer.CacheContainer {
	return &cachecontainer.CacheContainer{
		Any:    cache.New[any](rs),
		String: cache.New[string](rs),
	}
}
