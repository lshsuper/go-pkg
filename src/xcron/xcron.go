package xcron

import (
	"github.com/robfig/cron/v3"
	"sync"
)

//xCron 调度核心模型
type xCron struct {
	cron   *cron.Cron
	lock   *sync.RWMutex
	jobMap map[string]int
}

//IJob job接口
type IJob interface {
	cron.Job
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

//Start 开启
func (c *xCron) Start() {

	c.cron.Start()

}

//Put 塞入job
func (c *xCron) Put(job IJob) (id int, err error) {
	if len(job.JsonName()) <= 0 {
		panic("请实现JobName()方法...")
	}

	defer c.lock.Unlock()
	c.lock.Lock()
	i, err := c.cron.AddJob(job.Cron(), job)
	id = int(i)

	if err != nil {
		return 0, err
	}
	//记录下来
	c.jobMap[job.JsonName()] = int(id)
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

//Add 添加函数job
func (c *xCron) Add(cron, jobName string, fn func()) (id int, err error) {

	defer c.lock.Unlock()
	c.lock.Lock()

	i, err := c.cron.AddFunc(cron, fn)
	id = int(i)

	if err != nil {
		return
	}

	//记录下来
	c.jobMap[jobName] = int(id)
	return

}

//Stop 停止调度
func (c *xCron) Stop() {
	c.jobMap = nil
	c.cron.Stop()
}
