package cache

import (
	"github.com/garyburd/redigo/redis"
	"go-redis/pkg/setting"
	"time"
)

func NewRedisEngine(cacheSetting *setting.CacheSettingS) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     cacheSetting.MaxIdle,
		MaxActive:   cacheSetting.MaxActive,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", cacheSetting.Host)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}, nil
}
