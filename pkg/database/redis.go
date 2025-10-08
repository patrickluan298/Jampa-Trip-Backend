package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient - cliente Redis global
var RedisClient *redis.Client

// RedisConfig - configuração do Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

// RedisClientNew - inicializa o cliente Redis
func RedisClientNew() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}

	db, err := strconv.Atoi(Config.RedisDB)
	if err != nil {
		log.Fatalf("Erro ao converter REDIS_DB para int: %s", err.Error())
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Config.RedisHost, Config.RedisPort),
		Password: Config.RedisPassword,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar com Redis: %s", err.Error())
	}

	return RedisClient
}

// RedisPing - verifica se o Redis está funcionando
func RedisPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	return err
}
