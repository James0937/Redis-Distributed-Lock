package worker

import (
	"fmt"
	"log"
	"time"
)

const LOCK = "redislock"

func Executor() {
	locker := NewRedisLock(LOCK)
	err := locker.TryLock()
	if err != nil {
		log.Println(err)
		return
	}
	defer locker.UnLock()

	fmt.Println("hello")
	time.Sleep(5 * time.Second)
	fmt.Println("finish")
}
