package xredis

import (
	"go-pkg/src/xredis/domain"
	"sync"
)

//RedisFatory 全局实例
var RedisFatory *redisFactory

type redisFactory struct {
	lock     *sync.RWMutex
	redisMap map[string]*xRedis
}

//Register 注册Redis操作工厂
func Register() {
	RedisFatory = &redisFactory{
		lock:     new(sync.RWMutex),
		redisMap: map[string]*xRedis{},
	}
}

//Set 注册
func (f *redisFactory) Set(opt domain.RedisOpt) *redisFactory {
	defer f.lock.Unlock()
	f.lock.Lock()

	o := NewRedis(opt)
	f.redisMap[opt.Tag] = o
	return f
}

//Client redis-client
func (f *redisFactory) Client(tag string) *xRedis {
	defer f.lock.RUnlock()
	f.lock.RLock()
	o, _ := f.redisMap[tag]
	return o
}
