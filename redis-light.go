package redislight // Where is an Entry point?

import "errors"

// ErrKeyIsNotExists indicates that there is no value associated with provided key
var ErrKeyIsNotExists = errors.New("key is not exists")

// Storage provides abstraction for data storage
type Storage interface {
	Get(key string) (string, error)
	Set(key, val string) error
	Del(key string) error
}
