package producer

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheClient struct {
	client *memcache.Client
}

func NewMemcacheClient(address string) *MemcacheClient {
	return &MemcacheClient{
		client: memcache.New(address),
	}
}

func (c *MemcacheClient) Set(key string, value []byte, expiration time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(expiration.Seconds()),
	}

	err := c.client.Set(item)
	if err != nil {
		log.Println("Error setting value in Memcache:", err)
		return err
	}

	return nil
}

func (c *MemcacheClient) Get(key string) ([]byte, error) {
	item, err := c.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		log.Println("Error getting value from Memcache:", err)
		return nil, err
	}

	return item.Value, nil
}

func (c *MemcacheClient) Delete(key string) error {
	err := c.client.Delete(key)
	if err != nil {
		log.Println("Error deleting value from Memcache:", err)
		return err
	}
	return nil
}
