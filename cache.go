package main

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	rdb *redis.Client
}

func createCache(size int) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Cache{rdb: rdb}, nil
}

func (c *Cache) cacheURL(shortURL string, URL string) {
	_, err := (*c.rdb).Get(ctx, shortURL).Result()
	if err == redis.Nil {
		(*c.rdb).Set(ctx, shortURL, URL, 0)
	}
}

func (c *Cache) getURL(shortURL string) (string, error) {
	url, err := (*c.rdb).Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", errors.New("URL not found in cache")
	}
	return url, nil
}
