package cache

import "time"

type Cacher interface {
	// Origin get the underlying cache implementation.
	Origin() interface{}
	// Get Retrieve an item from the cache by key.
	Get(key string) (string, error)
	// Set a value to cache
	Set(key, val string, dur time.Duration) error
	// Forever Store an item in the cache indefinitely.
	Forever(key, val string) error
	// Forget Remove an item from the cache.
	Forget(key string) error
	//Increment the value of an item in the cache.
	Increment(key string, step ...int64) (int64, error)
	//Decrement the value of an item in the cache.
	Decrement(key string, step ...int64) (int64, error)
}
