package config

import (
	"os"
	"strconv"
	"time"
)

type MongoConfig struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
}

type RedisConfig struct {
	Addr            string
	Password        string
	DB              int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	MinIdleConns    int
	ConnMaxIdleTime time.Duration
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type Config struct {
	Mongo MongoConfig
	Redis RedisConfig
	JWT   JWTConfig
}

func Load() Config {
	return Config{
		Mongo: MongoConfig{
			URI:            getEnv("MONGO_URI", "mongodb://mongo:mongo@localhost:27017/?authSource=admin"),
			Database:       getEnv("MONGO_DB", "gps"),
			ConnectTimeout: getEnvDurationSeconds("MONGO_CONNECT_TIMEOUT_SECONDS", 10),
		},
		Redis: RedisConfig{
			Addr:            getEnv("REDIS_ADDR", "localhost:6379"),
			Password:        getEnv("REDIS_PASSWORD", ""),
			DB:              getEnvInt("REDIS_DB", 0),
			DialTimeout:     getEnvDurationSeconds("REDIS_DIAL_TIMEOUT_SECONDS", 5),
			ReadTimeout:     getEnvDurationSeconds("REDIS_READ_TIMEOUT_SECONDS", 3),
			WriteTimeout:    getEnvDurationSeconds("REDIS_WRITE_TIMEOUT_SECONDS", 3),
			PoolSize:        getEnvInt("REDIS_POOL_SIZE", 10),
			MinIdleConns:    getEnvInt("REDIS_MIN_IDLE_CONNS", 2),
			ConnMaxIdleTime: getEnvDurationSeconds("REDIS_CONN_MAX_IDLE_SECONDS", 300),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "super-secret-key"),
			Expiry: getEnvDuration("JWT_EXPIRY", time.Hour),
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvDurationSeconds(key string, fallbackSeconds int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallbackSeconds) * time.Second
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(fallbackSeconds) * time.Second
	}
	return time.Duration(parsed) * time.Second
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
