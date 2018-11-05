package common

import (
	"fmt"
	"sync"
)


type Process func(Job)
type Job struct {
	Workload interface{}
}
func NewJob(load interface{}) Job {
	return Job{Workload: load}
}

type Worker struct {
	WorkerPool chan chan Job
	Jobs chan Job //worker的工作队列
	quit chan bool //是否停止接单
}
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool:workerPool,
		Jobs: make(chan Job),
		quit: make(chan bool),
	}
}

/* worker 开工 */
func (w Worker) Start(do Process) {
	go func() {
		for {
			w.WorkerPool <- w.Jobs  //通知master，现在空闲
			select {
			case job := <-w.Jobs:
				do(job)
			case <-w.quit:
				return
			}
		}
	}()
}


type Master struct {
	Workers chan chan Job //worker池，实际上worker是一个job链
	jobQueue chan Job //待处理的任务chan
	num int  // worker 的个数
}


var Dispather *Master
var once sync.Once
//maxWorkers:开启线程数
func NewMaster(maxWorkers int) *Master {
	return &Master{Workers: make(chan chan Job, maxWorkers),  jobQueue: make(chan Job, 2*maxWorkers), num:maxWorkers}
}

func GetInstance(maxWorkers int) *Master {

	once.Do(func() {
		Dispather = NewMaster(maxWorkers)
	})
	return Dispather
}

func (m *Master) Run(do Process) {
	//启动所有的Worker
	for i := 0; i < m.num; i++ {
		work := NewWorker(m.Workers)
		work.Start(do)
	}
	go m.dispatch()
}
func (m *Master) dispatch() {
	for {
		select {
		case job := <-m.jobQueue:
			go func(job Job) {
				//从Workers中取出一个worker
				worker := <-m.Workers
				//向这个worker派发job
				worker <- job
			}(job)
		}
	}
}
//添加任务到任务通道
func (m *Master) AddJob(load interface{}) {
	job := NewJob(load)
	//向任务通道发送任务
	m.jobQueue <- job
	fmt.Printf("job pool len is %d \n", len(m.jobQueue))
}

