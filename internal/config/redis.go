package config

import "user-management-api/internal/utils"

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRedisConfig() {
	cfg := RedisConfig{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Username: utils.GetEnv("REDIS_USER", ""),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetIntEnv("REDIS_DB", 0),
	}
}
