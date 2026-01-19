package deps

import (
	"context"
	"gps/internal/adapters/repo/mongoDb"
	redisRepo "gps/internal/adapters/repo/redis"
	"gps/internal/app_services/aggregator"
	"gps/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Deps struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
	MongoRepo       *mongoDb.Repository
	Redis       *redisRepo.Repository
	Aggregator *aggregator.AggregatorService
	
}
type option func(*Deps) error

func NewDeps(opts ...option) (*Deps, error) {
	deps := &Deps{}
	for _, opt := range opts {
		if err := opt(deps); err != nil {
			return nil, err
		}
	}
	return deps, nil
}
func WithMongoRepo(config config.Config) option {
	return func(d *Deps) error {
		repo, err := mongoDb.NewRepository(context.Background(), d.MongoClient, config.Mongo.Database)
		if err != nil {
			return err
		}
		d.MongoRepo = repo
		return nil
	}
}
func WithRedisRepo(config config.Config) option {
	return func(d *Deps) error {
		repo, err := redisRepo.NewRepository(d.RedisClient, 5*time.Minute)
		if err != nil {
			return err
		}
		d.Redis = repo
		return nil
	}
}

func WithMongoClient(ctx context.Context, config config.Config) option {
	return func(d *Deps) error {
		client, err := mongo.Connect(options.Client().ApplyURI(config.Mongo.URI))
		if err != nil {
			return err
		}

		pingCtx, cancel := context.WithTimeout(ctx, config.Mongo.ConnectTimeout)
		defer cancel()
		if err := client.Ping(pingCtx, nil); err != nil {
			_ = client.Disconnect(context.Background())
			return err
		}
		d.MongoClient = client
		return nil
	}
}

func WithRedisClient(ctx context.Context, config config.Config) option {
	return func(d *Deps) error {
		client := redis.NewClient(&redis.Options{
			Addr:            config.Redis.Addr,
			Password:        config.Redis.Password,
			DB:              config.Redis.DB,
			DialTimeout:     config.Redis.DialTimeout,
			ReadTimeout:     config.Redis.ReadTimeout,
			WriteTimeout:    config.Redis.WriteTimeout,
			PoolSize:        config.Redis.PoolSize,
			MinIdleConns:    config.Redis.MinIdleConns,
			ConnMaxIdleTime: config.Redis.ConnMaxIdleTime,
		})

		if err := client.Ping(ctx).Err(); err != nil {
			_ = client.Close()
			return err
		}
		d.RedisClient = client
		return nil
	}
}
