// Package container is a package that provides various containers for the application.
package container

import "github.com/eko/gocache/lib/v4/cache"

// NewCache is a function that returns a new container.Cache instance.
func NewCache(a *cache.Cache[any], s *cache.Cache[string]) *Cache {
	return &Cache{
		Any:    a,
		String: s,
	}
}

// Cache is a struct that contains all of the cache.Cache instances.
type Cache struct {
	// Any is a cache.Cache instance that is used for storing any type of data.
	Any *cache.Cache[any]
	// String is a cache.Cache instance that is used for storing strings.
	String *cache.Cache[string]
}
