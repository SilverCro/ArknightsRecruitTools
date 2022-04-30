package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"kratosTest/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewRedis, NewGreeterRepo, NewMyKratosRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	rdb *redis.Client
}

//// NewData .
//func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
//	cleanup := func() {
//		log.NewHelper(logger).Info("closing the data resources")
//	}
//	return &Data{}, cleanup, nil
//}

func NewRedis(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log.NewHelper(logger).Info("redisAddr", c.Redis.Addr)
	log.NewHelper(logger).Info("redisPassword", c.Redis.Password)
	log.NewHelper(logger).Info("redisDB", c.Redis.DB)
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password, // no password set
		DB:           int(c.Redis.DB),  // use default DB
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	cleanup := func() {
		err := rdb.Close()
		if err != nil {
			return
		}
	}
	return &Data{rdb: rdb}, cleanup, nil
}
