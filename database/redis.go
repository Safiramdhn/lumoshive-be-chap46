package database

import (
	"context"
	"fmt"
	"golang-chap46/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher struct {
	rdb      *redis.Client
	expiracy time.Duration
	prefix   string
}

func newRedisClient(url, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbIndex,
	})
}
func NewCacher(cfg config.Configuration, expiracy int) Cacher {
	client := newRedisClient(cfg.RedisConfig.Url, cfg.RedisConfig.Password, 0)

	// Test Redis connection and log the result
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connected successfully!")
	}

	cache := Cacher{
		rdb:      client,
		expiracy: time.Duration(expiracy) * time.Second,
		prefix:   cfg.RedisConfig.Prefix,
	}

	return cache
}

func (c *Cacher) Push(name string, value []byte) error {
	return c.rdb.RPush(context.Background(), c.prefix+"_"+name, value).Err()
}

func (c *Cacher) Pop(name string) (string, error) {
	return c.rdb.LPop(context.Background(), c.prefix+"_"+name).Result()
}

func (c *Cacher) GetLength(name string) int64 {
	return c.rdb.LLen(context.Background(), c.prefix+"_"+name).Val()
}

func (c *Cacher) Set(name string, value string) error {
	return c.rdb.Set(context.Background(), c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *Cacher) SaveToken(name string, value string) error {
	log.Printf("Saving token for %s:%s", name, value)
	return c.rdb.Set(context.Background(), c.prefix+"_"+name, value, 20*time.Hour).Err()
}

func (c *Cacher) Get(name string) (string, error) {
	result, err := c.rdb.Get(context.Background(), c.prefix+"_"+name).Result()
	if err != nil {
		fmt.Println("error getting token from Redis: ", err)
	}
	return result, err
}

func (c *Cacher) Delete(name string) error {
	return c.rdb.Del(context.Background(), c.prefix+"_"+name).Err()
}

func (c *Cacher) DeleteByKey(key string) error {
	return c.rdb.Del(context.Background(), key).Err()
}

func (c *Cacher) PrintKeys() {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println("key", key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}
}

func (c *Cacher) GetKeys() []string {
	var cursor uint64
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		result = append(result, keys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	return result
}

func (c *Cacher) GetKeysByPattern(pattern string) []string {
	var cursor uint64
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, pattern, 0).Result()
		if err != nil {
			panic(err)
		}

		result = append(result, keys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	return result
}

// Pub and Sub
func (c *Cacher) Publish(channelName string, message string) error {
	return c.rdb.Publish(context.Background(), channelName, message).Err()
}

func (c *Cacher) Subcribe(channelName string) (*redis.Message, error) {
	subscriber := c.rdb.Subscribe(context.Background(), channelName)
	message, err := subscriber.ReceiveMessage(context.Background())
	return message, err
}

func (c *Cacher) GetRedisClient() *redis.Client {
	return c.rdb
}
