package xredis

import (
	"context"
	"github.com/lshsuper/go-pkg/src/xredis/domain"
	"testing"
	"time"
)

func TestLock(t *testing.T) {

	r := NewRedis(domain.RedisOpt{
		Pwd:  "honghe@2020",
		Addr: ":6379",
		Tag:  "abc",
	})

	for i := 0; i < 100; i++ {

		go func() {
			r.LockAndTake(context.Background(), "orders", time.Second*time.Duration(5), func() {
				time.Sleep(time.Second * 1)
			})
		}()

		time.Sleep(time.Second * 1)

	}
	time.Sleep(time.Second * 100)

}
