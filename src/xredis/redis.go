package xredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lshsuper/go-pkg/src/xredis/domain"
	"sync"
	"time"
)

type xRedis struct {
	xClient *redis.Client
	lock    *sync.RWMutex
}

//NewRedis 注册Redis
func NewRedis(opt domain.RedisOpt) *xRedis {

	client := redis.NewClient(&redis.Options{
		Addr:     opt.Addr, // use default Addr
		Password: opt.Pwd,  // no password set
		DB:       0,        // use default DB
		PoolSize: opt.PoolSize,
	})

	cmd := client.Ping(context.Background())

	if err := cmd.Err(); err != nil {
		panic(fmt.Sprintf("redis[%s]初始化失败:%v", opt.Tag, err))
	}

	xClient := new(xRedis)
	xClient.xClient = client
	xClient.lock = new(sync.RWMutex)
	return xClient
}

//Publish 生产端
func (c *xRedis) Publish(ctx context.Context, channel, message string) error {
	cmd := c.xClient.Publish(ctx, channel, message)
	return cmd.Err()
}

//Consumer 消费端
func (c *xRedis) Consumer(ctx context.Context, channel string, callback func(msg string)) {

	go func() {
		pub := c.xClient.Subscribe(ctx, channel)

		for {
			msg, _ := pub.ReceiveMessage(ctx)
			callback(msg.Payload)
		}

	}()

}

//StringSet 设置字符串类型数据
func (c *xRedis) StringSet(ctx context.Context, key, value string, expiration time.Duration) error {
	cmd := c.xClient.Set(ctx, key, value, expiration)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}
	return nil
}

//StringGet 获取字符串类型数据
func (c *xRedis) StringGet(ctx context.Context, key string) (string, error) {
	cmd := c.xClient.Get(ctx, key)
	str, err := cmd.Result()
	if err != nil {
		//TODO:logger
		return "", err
	}
	return str, err
}

//HashMSet hash-set(支持批量键值)
func (c *xRedis) HashMSet(ctx context.Context, key string, value map[string]interface{}) error {
	cmd := c.xClient.HMSet(ctx, key, value)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}

	return nil
}

//HashMGet hash-get(支持批量获取)
//notice:当对应的field不存在时,对用map的value为nil
func (c *xRedis) HashMGet(ctx context.Context, key string, fields ...string) (map[string]interface{}, error) {
	cmd := c.xClient.HMGet(ctx, key, fields...)
	res, err := cmd.Result()
	if err != nil {
		//TODO:logger
		return nil, err
	}
	//组装map
	resMap := make(map[string]interface{}, len(res))
	for k, v := range res {
		resMap[fields[k]] = v
	}

	return resMap, nil
}

//HashDel 删除指定key的hash下的指定fields
func (c *xRedis) HashDel(ctx context.Context, key string, fields string) error {
	cmd := c.xClient.HDel(ctx, key, fields)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}
	return nil
}

//Del 产出指定key的缓存
func (c *xRedis) Del(ctx context.Context, keys ...string) error {

	cmd := c.xClient.Del(ctx, keys...)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}
	return nil

}

//BitMapSet  设置
func (c *xRedis) BitMapSet(ctx context.Context, key string, offset, value int) error {

	cmd := c.xClient.SetBit(ctx, key, int64(offset), value)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}
	return nil

}

//BitMapCount  获取
func (c *xRedis) BitMapCount(ctx context.Context, key string) (int, error) {

	cmd := c.xClient.BitCount(ctx, key, nil)
	res, err := cmd.Result()
	if err != nil {
		//TODO:logger
		return 0, err
	}
	return int(res), err

}

//ListPush 列表入参
func (c *xRedis) ListPush(ctx context.Context, p domain.Position, key string, value ...interface{}) error {

	//左操作
	if p == domain.Left {
		cmd := c.xClient.LPush(ctx, key, value...)
		if _, err := cmd.Result(); err != nil {
			//TODO:logger
			return err
		}
		return nil
	}

	cmd := c.xClient.RPush(ctx, key, value...)
	if _, err := cmd.Result(); err != nil {
		//TODO:logger
		return err
	}
	return nil

}

//ListPop 出队
func (c *xRedis) ListPop(ctx context.Context, p domain.Position, key string) (interface{}, error) {

	//左操作
	if p == domain.Left {
		cmd := c.xClient.LPop(ctx, key)
		res, err := cmd.Result()
		if err != nil {
			//TODO:logger
			return nil, err
		}
		return res, err
	}

	cmd := c.xClient.RPop(ctx, key)
	res, err := cmd.Result()
	if err != nil {
		//TODO:logger
		return nil, err
	}
	return res, err

}

//LockAndTake 锁并执行
func (c *xRedis) LockAndTake(ctx context.Context, key string, expiration time.Duration, callback func()) {

	cmd := c.xClient.SetNX(ctx, key, 1, expiration)
	res, err := cmd.Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !res {
		return
	}

	callback()
	c.xClient.Del(ctx, key)
	return

}

func (c *xRedis) OptimisticLock(ctx context.Context, key string, expiration time.Duration) error {

	return c.xClient.Watch(ctx, func(tx *redis.Tx) error {
		// 操作仅在 Watch 的 Key 没发生变化的情况下提交
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, "_lock", expiration)
			return nil
		})
		return err
	}, key)

}
//HashLen hash列表数量
func (c *xRedis)HashLen(ctx context.Context,key string)(int,error)  {

	cmd := c.xClient.HLen(ctx,key)
	num,err:=cmd.Result()
    return int(num),err
}
