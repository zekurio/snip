package database

import (
	"strings"
	"sync"
)

const (
	cacheKeySeparator = ":"
)

type DatabaseCache struct {
	IDatabase

	cache sync.Map
}

var _ IDatabase = (*DatabaseCache)(nil)

func WrapCache(inheritedDB IDatabase, inheritedErr error) (wrappedDB IDatabase, err error) {
	if err != nil {
		return nil, err
	}

	var cache DatabaseCache
	cache.IDatabase = inheritedDB

	return &cache, nil
}

func ckey(elements ...string) string {
	return strings.Join(elements, cacheKeySeparator)
}
