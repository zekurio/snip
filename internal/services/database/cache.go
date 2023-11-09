package database

import (
	"strings"
	"sync"

	"github.com/zekurio/snip/internal/models"
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

func (c *DatabaseCache) AddUpdateLink(link *models.Link) error {
	err := c.IDatabase.AddUpdateLink(link)
	if err != nil {
		return err
	}

	c.cache.Store(ckey("link", link.ID), link)

	return nil
}

func (c *DatabaseCache) GetLinkByID(id string) (*models.Link, error) {
	key := ckey("link", id)

	if val, ok := c.cache.Load(key); ok {
		return val.(*models.Link), nil
	}

	link, err := c.IDatabase.GetLinkByID(id)
	if err != nil {
		return nil, err
	}

	c.cache.Store(key, link)

	return link, nil
}

func (c *DatabaseCache) DeleteLink(id string) error {
	err := c.IDatabase.DeleteLink(id)
	if err != nil {
		return err
	}

	c.cache.Delete(ckey("link", id))

	return nil
}

func ckey(elements ...string) string {
	return strings.Join(elements, cacheKeySeparator)
}
