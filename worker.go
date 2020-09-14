package plugins

import (
	"os"
	"runtime"
	"strconv"
)

var (
	defaultMaxQueue = 20000
)

// Job 表示要运行的Job的接口
type Job interface {
	Do()
}

// JobQueue Job 通道
var JobQueue chan Job

// Worker 表示执行该Job的worker
type Worker struct {
	// WorkerPool 是个指向全局唯一的 chan 的引用,
	// 负责传递 Worker 接收 Job 的 chan。
	// Worker 空闲时，将自己的 JobChannel 放入 WorkerPool 中。
	// Dispatcher 收到新的 Job 时，从 JobChannel 中取出一个 chan， 并将 Job
	// 放入其中，此时 Worker 将从 Chan 中接收到 Job，并进行处理
	WorkerPool chan chan Job
	// Worker 用于接收 Job 的 chan
	JobChannel chan Job
	// 用于给 Worker 发送控制命令的 chan，用于停止 chan
	Quit chan bool
}

func NewWorker(workPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
	}
}

// Start 为worker启动运行循环，侦听退出通道
func (w Worker) Start() {
	go func() {
		for {
			// 将当前worker注册到worker队列中
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// 激活 Job
				job.Do()
			case <-w.Quit:
				// 收到停止信号
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// Dispatcher 调度器
type Dispatcher struct {
	// 工作池容量
	MaxWorkers int
	// 向调度程序注册的工作程序通道池
	WorkerPool chan chan Job
	Quit       chan bool
}

// NewDispatcher 新建一个调度器
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{MaxWorkers: maxWorkers, WorkerPool: pool, Quit: make(chan bool)}
}

// Run 运行调度器
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.Dispatch()
}

// Stop 调度器停止
func (d *Dispatcher) Stop() {
	go func() {
		d.Quit <- true
	}()
}

// Dispatch 分配任务
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// 已收到Job请求
			go func(job Job) {
				// 尝试获取可用的worker Job通道
				// 阻塞直到worker空闲为止
				jobChannel := <-d.WorkerPool
				// 将Job分配到worker
				jobChannel <- job
			}(job)
		case <-d.Quit:
			return
		}
	}
}

/*func init() {
	maxQueue, err := strconv.Atoi(os.Getenv("MAX_JOB_QUEUE"))
	if err != nil || maxQueue <= 0 {
		maxQueue = defaultMaxQueue
	}
	JobQueue = make(chan Job, maxQueue)
	dispatcher := NewDispatcher(runtime.NumCPU())
	dispatcher.Run()
}*/
