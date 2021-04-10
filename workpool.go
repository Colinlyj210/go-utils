package go_utils

import "log"

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

func (workPool *WorkPool) startWork(id int) {
	for {
		select {
		case task := <-workPool.taskQueue:
			log.Print("work id: ", id)
			task.run(workPool)
		}
	}
}

func (workPool *WorkPool) ExecuteTask(t *Task) {
	workPool.taskQueue <- t
}

func (workPool *WorkPool) Execute(f TaskFunc, args ...interface{}) {
	workPool.taskQueue <- &Task{f, args}
}

func NewWorkPool(workerNum, taskQueueSize int) *WorkPool {
	workPool := &WorkPool{
		taskQueue: make(chan *Task, taskQueueSize),
	}
	for i := 0; i < workerNum; i++ {
		go workPool.startWork(i)
	}
	return workPool
}
