package provider

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCacheProvider() *cache.Cache {
	c := cache.New(5*time.Second, 1*time.Minute)

	return c
}
