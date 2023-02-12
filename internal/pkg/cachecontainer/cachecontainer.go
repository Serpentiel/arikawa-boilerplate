// Package cachecontainer is the package that contains the cache container.
package cachecontainer

import "github.com/eko/gocache/lib/v4/cache"

// CacheContainer is a struct that contains all of the cache.Cache instances.
type CacheContainer struct {
	// Any is a cache.Cache instance that is used for storing any type of data.
	Any *cache.Cache[any]
	// String is a cache.Cache instance that is used for storing strings.
	String *cache.Cache[string]
}
