package executor

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrStoppedWorker = fmt.Errorf("worker is stopped")
	ErrStartedWorker = fmt.Errorf("worker is already started")
)

const (
	WORK    = 1
	STOPPED = 0
)

type AsyncExecutor struct {
	queue 	     *chan func()
	wg    		 *sync.WaitGroup
	workersCount int
	isWork       *int32
}

func NewAsyncExecutor(workersCount int) AsyncExecutor {
	var isWork int32 = 0
	worker := AsyncExecutor{
		wg:    		  &sync.WaitGroup{},
		workersCount: workersCount,
		isWork:       &isWork, 
	}

	tmp := make(chan func(), workersCount)
	worker.queue = &tmp

	return worker
}

func (instance *AsyncExecutor) Start() error {
	if !atomic.CompareAndSwapInt32(instance.isWork, STOPPED, WORK) {
		return ErrStartedWorker
	}

	*instance.queue = make(chan func(), instance.workersCount)

	for i := 0; i < instance.workersCount; i++ {
		go func() {
			for task := range *instance.queue {
				task()
				instance.wg.Done()
			}
		}()
	}

	return nil
}

func (instance *AsyncExecutor) Stop() error {
	if atomic.CompareAndSwapInt32(instance.isWork, WORK, STOPPED) {
		return ErrStoppedWorker
	}

	close(*instance.queue)

	return nil
}

func (instance *AsyncExecutor) AddTask(task func()) error {
	if !atomic.CompareAndSwapInt32(instance.isWork, WORK, WORK) {
		return ErrStoppedWorker
	}

	instance.wg.Add(1)
	*instance.queue <- task

	return nil
}

func (instance *AsyncExecutor) Wait() error {
	if !atomic.CompareAndSwapInt32(instance.isWork, WORK, WORK){
		return ErrStoppedWorker
	}

	instance.wg.Wait()

	return nil
}