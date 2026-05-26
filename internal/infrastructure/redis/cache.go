// infrastructure/redis/cache.go
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// Guardar en caché
func (r *CacheRepository) Set(key string, value []byte, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

// Leer de caché
func (r *CacheRepository) Get(key string) ([]byte, error) {
	result, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Eliminar de caché
func (r *CacheRepository) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
