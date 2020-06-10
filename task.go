package tcc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mushroomsir/tcc/store"
)

// Task ...
type Task struct {
	Name      string
	Value     string
	CreatedAt time.Time

	uid    string
	store  store.TccAsyncTaskInterface
	logger LoggerInterface
}

// Confirm submit async task
func (a *Task) Confirm() error {
	err := a.store.Confirm(a.uid)
	if err != nil {
		a.logger.Err(fmt.Errorf("Confirm error, %s", err.Error()))
	}
	return err
}

// Cancel ...
func (a *Task) Cancel() error {
	err := a.store.Cancel(a.uid)
	if err != nil {
		a.logger.Err(fmt.Errorf("cancel error, %s", err.Error()))
	}
	return err
}

// JSONToObj ...
func (a *Task) JSONToObj(dest interface{}) error {
	err := json.Unmarshal([]byte(a.Value), dest)
	return err
}
