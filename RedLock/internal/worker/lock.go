package worker

import (
	"context"
	"errors"
	"github.com/garyburd/redigo/redis"
	"go-redis/global"
	"log"
	"time"
)

type RedisLock struct {
	Key        string
	TTL        int64
	IsLocked   bool
	CancelFunc context.CancelFunc
}

func NewRedisLock(key string) *RedisLock {
	redisLock := &RedisLock{
		Key: key,
		TTL: 3,
	}
	return redisLock
}

func (r *RedisLock) TryLock() (err error) {
	if err = r.Grant(); err != nil {
		return
	}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	r.CancelFunc = cancelFunc
	r.KeepAlive(ctx)
	r.IsLocked = true
	return nil
}

func (r *RedisLock) UnLock() (err error) {
	var res int
	if r.IsLocked {
		if res, err = redis.Int(global.GetConn().Do("DEL", r.Key)); err != nil {
			return errors.New("ReleaseFailure")
		}
		if res == 1 {
			r.CancelFunc()
			return
		}
	}
	return errors.New("ReleaseFailure")
}

func (r *RedisLock) KeepAlive(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := redis.Int(global.GetConn().Do("EXPIRE", r.Key, r.TTL))
				if err != nil {
					log.Println("AutoRelease", err)
				}
				if res != -1 {
					log.Println("AutoRenew")
				}
				time.Sleep(time.Duration(r.TTL/2) * time.Second)
			}
		}
	}()
}

func (r *RedisLock) Grant() (err error) {
	if res, err := redis.String(global.GetConn().Do("SET", r.Key, "1", "NX", "EX", r.TTL)); err == nil {
		if res == "OK" {
			return nil
		}
	}
	return errors.New("LockFailure")
}
