package main

import (
	"go-redis/global"
	"go-redis/internal/worker"
	"go-redis/pkg/cache"
	"go-redis/pkg/setting"
	"log"
	"strings"
	"sync"
)

func Init(config string) {
	err := setupSetting(config)
	if err != nil {
		log.Printf("init setupSetting err: %v\n", err)
	} else {
		log.Printf("Initialization Configuration Success")
	}
	err = setupCacheEngine()
	if err != nil {
		log.Printf("init setupCacheEngine err: %v\n", err)
	} else {
		log.Printf("Initialization Cache Success")
	}
}

func setupSetting(config string) error {
	newSetting, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("Cache", &global.CacheSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupCacheEngine() error {
	var err error
	global.RedisEngine, err = cache.NewRedisEngine(global.CacheSetting)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	Init("configs/")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.Executor()
		}()
	}
	wg.Wait()
}
