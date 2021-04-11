package go_utils

type TaskFunc func(wp *WorkPool, args ...interface{})

type Task struct {
	taskFunc TaskFunc
	args     []interface{}
}

func (task *Task) run(workPool *WorkPool) {
	task.taskFunc(workPool, task.args...)
}

type WorkPool struct {
	taskQueue chan *Task
}

func (workPool *WorkPool) startWork() {
	for {
		select {
		case task := <-workPool.taskQueue:
			task.run(workPool)
		}
	}
}

func (workPool *WorkPool) ExecuteTask(task *Task) {
	workPool.taskQueue <- task
}

func (workPool *WorkPool) Execute(taskFunc TaskFunc, args ...interface{}) {
	workPool.taskQueue <- &Task{taskFunc, args}
}

func NewWorkPool(workerNum, taskQueueSize int) *WorkPool {
	workPool := &WorkPool{
		taskQueue: make(chan *Task, taskQueueSize),
	}
	for i := 0; i < workerNum; i++ {
		go workPool.startWork()
	}
	return workPool
}
