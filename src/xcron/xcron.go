package xcron

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type xCron struct {
	cron   *cron.Cron
	lock   *sync.RWMutex
	jobMap map[string]int
}

type IJob interface {
	Run()
	JsonName() string
	Cron() string
}

func NewXCron() *xCron {

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithSeconds())
	return &xCron{
		cron:   c,
		lock:   new(sync.RWMutex),
		jobMap: make(map[string]int, 0),
	}

}

//Put 塞入job
func (c *xCron) Put(job IJob) (id cron.EntryID, err error) {
	if len(job.JsonName()) <= 0 {
		panic("请实现JobName()方法...")
	}

	defer c.lock.Unlock()
	c.lock.Lock()
	id, err = c.cron.AddJob(job.Cron(), job)

	if err != nil {
		//记录下来
		c.jobMap[job.JsonName()] = int(id)
	}
	return
}

//Remove 移除job
func (c *xCron) Remove(jobName string) {
	defer c.lock.Unlock()
	c.lock.Lock()
	j, ok := c.jobMap[jobName]
	if ok {
		c.cron.Remove(cron.EntryID(j))
	}
	return
}
