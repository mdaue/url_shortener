package main

import (
	"errors"
)

type Cache struct {
	cache *map[string]string
	size  int
}

func createCache(size int) (*Cache, error) {
	cache := make(map[string]string, size)
	return &Cache{cache: &cache, size: size}, nil
}

func (c *Cache) cacheURL(shortURL string, URL string) {
	if len(*c.cache) >= c.size {
		c.pruneCache()
	}
	_, ok := (*c.cache)[shortURL]
	if !ok {
		(*c.cache)[shortURL] = URL
	}
}

func (c *Cache) getURL(shortURL string) (string, error) {
	url, ok := (*c.cache)[shortURL]
	if !ok {
		return "", errors.New("URL not found in cache")
	}
	return url, nil
}

func (c *Cache) pruneCache() {
	count := 0
	target := len(*c.cache) / 2
	for shortURL := range *c.cache {
		delete(*c.cache, shortURL)
		count++
		if count >= target {
			break
		}
	}
}
