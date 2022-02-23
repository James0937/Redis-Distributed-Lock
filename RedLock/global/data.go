package global

import "github.com/garyburd/redigo/redis"

var (
	RedisEngine *redis.Pool
)

func GetConn() redis.Conn {
	return RedisEngine.Get()
}
