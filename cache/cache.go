package cache

import "time"

// Cache interface contains all behaviors for cache adapter.
type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// get a batch cached values.
	GetMulti(keys []string) []interface{}
	// set cached value.
	Put(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// check whether cached value exists.
	IsExist(key string) bool
	// increment counter.
	Increment(key string) error
	// decrement counter.
	Decrement(key string) error
	// clear all cache.
	ClearAll() error
}
