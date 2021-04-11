package go_utils

import (
	"context"
	"sync"
)

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

	stopCtx        context.Context
	stopCancelFunc context.CancelFunc
	wg             sync.WaitGroup
}

func (workPool *WorkPool) startWork() {
	for {
		select {
		case <-workPool.stopCtx.Done():
			workPool.wg.Done()
			return
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
	workPool.Start(workerNum)
	return workPool
}

func (workPool *WorkPool) Start(workerNum int) {
	workPool.wg.Add(workerNum)
	workPool.stopCtx, workPool.stopCancelFunc = context.WithCancel(context.Background())
	for i := 0; i < workerNum; i++ {
		go workPool.startWork()
	}
}

func (workPool *WorkPool) Stop() {
	workPool.stopCancelFunc()
	workPool.wg.Wait()
}
