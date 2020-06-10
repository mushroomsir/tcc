package tcc

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mushroomsir/tcc/lock"
	"github.com/mushroomsir/tcc/store"
	"github.com/mushroomsir/tcc/store/dto"
)

const (
	tccTry     = "try"
	tccConfirm = "confirm"
)

type handler func(task *Task)

// Option ...
type Option struct {
	TaskConcurrency   int
	PullTaskInterval  int
	PullTaskBatchSize int
	LockExpire        time.Duration
	TryFailed         time.Duration

	Store  store.TccAsyncTaskInterface
	Lock   lock.LockInterface
	Logger LoggerInterface
}

// TCC ...
type TCC struct {
	Option

	TryHandler     handler
	ConfirmHandler handler

	rwMutex sync.RWMutex
}

// New ...
func New(option *Option) *TCC {
	a := &TCC{
		Option: *option,
	}
	a.init()
	go a.Loop()
	return a
}

func (a *TCC) init() {
	if a.Store == nil {
		panic("store is nil")
	}
	if a.Lock == nil {
		panic("lock is nil")
	}
	if a.Logger == nil {
		a.Logger = &Logger{}
	}
	if a.PullTaskInterval < 1 {
		a.PullTaskInterval = 3
	}
	if a.TaskConcurrency < 1 {
		a.TaskConcurrency = 2
	}
	if a.PullTaskBatchSize < 1 {
		a.PullTaskBatchSize = 20
	}
	if a.LockExpire < 1 {
		a.LockExpire = 5 * time.Second
	}
	if a.TryFailed < 1 {
		a.TryFailed = 3 * time.Second
	}
}

// SetTryHandler ...
func (a *TCC) SetTryHandler(h handler) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()
	a.TryHandler = h
}

// getTryHandler ...
func (a *TCC) getTryHandler() handler {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()
	return a.TryHandler
}

// SetConfirmHandler ...
func (a *TCC) SetConfirmHandler(h handler) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()
	a.ConfirmHandler = h
}

// SetConfirmHandler ...
func (a *TCC) getConfirmHandler() handler {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()
	return a.ConfirmHandler
}

// NewTransaction ...
func (a *TCC) NewTransaction(name string) *Transaction {
	t := &Transaction{
		uid:   uuid.New().String(),
		name:  name,
		store: a.Store,
	}
	return t
}

// Loop ...
func (a *TCC) Loop() {
	s := time.Duration(a.PullTaskInterval) * time.Second
	for {
		time.Sleep(s)
		tasks, err := a.Store.Tasks(a.PullTaskBatchSize)
		if err != nil {
			a.Logger.Err(fmt.Errorf("pull task error, %s", err.Error()))
			continue
		}
		if len(tasks) == 0 {
			continue
		}
		taskChannel := make(chan *dto.TccAsyncTaskSchema)

		var wg sync.WaitGroup
		wg.Add(a.TaskConcurrency)

		for i := 0; i < a.TaskConcurrency; i++ {
			go func() {
				defer wg.Done()
				for {
					task := <-taskChannel
					if task == nil {
						return
					}
					a.handler(task)
				}
			}()
		}
		for _, item := range tasks {
			if item.Status == tccTry && time.Now().Sub(item.CreatedAt) < a.TryFailed {
				continue
			}
			taskChannel <- item
		}
		close(taskChannel)
		wg.Wait()
	}
}

func (a *TCC) handler(task *dto.TccAsyncTaskSchema) {
	defer func() {
		if err := recover(); err != nil {
			a.Logger.Err(fmt.Errorf("handler task error, %v", err))
		}
	}()
	out := &Task{
		uid:       task.UID,
		store:     a.Store,
		Name:      task.Name,
		Value:     task.Value,
		CreatedAt: task.CreatedAt,
	}
	if task.Status == tccTry && a.getTryHandler() != nil {
		a.handlerTry(task, out)
	} else if task.Status == tccConfirm && a.getConfirmHandler() != nil {
		a.handlerConfirm(task, out)
	} else {
		err := out.Cancel()
		if err != nil {
			a.Logger.Err(fmt.Errorf("cancel error, %s", err.Error()))
		}
	}
}

func (a *TCC) handlerTry(task *dto.TccAsyncTaskSchema, out *Task) {
	err := a.Lock.Lock(task.UID, a.LockExpire)
	if err != nil {
		return
	}
	res, err := a.Store.Task(task.UID)
	if err != nil {
		a.Logger.Err(fmt.Errorf("get task error, %s", err.Error()))
	} else if res.Status == tccTry {
		a.getTryHandler()(out)
	}
	err = a.Lock.Unlock(task.UID)
	if err != nil {
		a.Logger.Err(fmt.Errorf("unlock error, %s", err.Error()))
	}
}

func (a *TCC) handlerConfirm(task *dto.TccAsyncTaskSchema, out *Task) {
	err := a.Lock.Lock(task.UID, a.LockExpire)
	if err != nil {
		return
	}
	res, err := a.Store.Task(task.UID)
	if err != nil {
		a.Logger.Err(fmt.Errorf("get task error, %s", err.Error()))
	} else if res.Status == tccConfirm {
		a.getConfirmHandler()(out)
	}
	err = a.Lock.Unlock(task.UID)
	if err != nil {
		a.Logger.Err(fmt.Errorf("unlock error, %s", err.Error()))
	}
}
