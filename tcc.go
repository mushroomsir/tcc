package tcc

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mushroomsir/logger/alog"
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

	Store store.TccAsyncTaskInterface
	Lock  lock.LockInterface
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
			alog.Err("pull task", err.Error())
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
			alog.Err("tcchandler", err)
		}
	}()
	t := &Task{
		uid:       task.UID,
		store:     a.Store,
		Name:      task.Name,
		Value:     task.Value,
		CreatedAt: task.CreatedAt,
	}
	if task.Status == tccTry && a.getTryHandler() != nil {
		err := a.Lock.Lock(task.UID, a.LockExpire)
		if err == nil {
			res, err := a.Store.Task(task.UID)
			if err == nil && res.Status == tccTry {
				a.getTryHandler()(t)
			}
			a.Lock.Unlock(task.UID)
		}
	} else if task.Status == tccConfirm && a.getConfirmHandler() != nil {
		err := a.Lock.Lock(task.UID, a.LockExpire)
		if err == nil {
			_, err = a.Store.Task(task.UID)
			if err == nil {
				a.getConfirmHandler()(t)
			}
			a.Lock.Unlock(task.UID)
		}
	} else {
		t.Cancel()
	}
}
