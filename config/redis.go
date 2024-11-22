package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func GetRedis() *redis.Client {
	cnf := GetConfig().RedisDB

	rdb := redis.NewClient(&redis.Options{
		Addr:     cnf.Addr,
		Password: cnf.Password,
		DB:       cnf.DB,
		Protocol: cnf.Protocol,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed connec redis: %v", err)
	}

	log.Println("connected redis response: " + result)
	return rdb

}
