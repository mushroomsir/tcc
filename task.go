package tcc

import (
	"encoding/json"
	"time"

	"github.com/mushroomsir/tcc/store"
)

// Task ...
type Task struct {
	Name      string
	Value     string
	CreatedAt time.Time

	uid   string
	store store.TccAsyncTaskInterface
}

// Confirm submit async task
func (a *Task) Confirm() error {
	return a.store.Confirm(a.uid)
}

// Cancel ...
func (a *Task) Cancel() error {
	return a.store.Cancel(a.uid)
}

// JsonToObj ...
func (a *Task) JsonToObj(dest interface{}) error {
	err := json.Unmarshal([]byte(a.Value), dest)
	return err
}
