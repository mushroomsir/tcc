package store

import (
	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/tcc/store/dto"
	"github.com/mushroomsir/tcc/store/mysql"
)

// TccAsyncTaskInterface ...
type TccAsyncTaskInterface interface {
	// TryPlan get try execute play, example sql.
	TryPlan(uid, name, value string) string
	// Try add async task
	Try(uid, name, value string) error
	// TryAndConfirm add and submit async task
	TryAndConfirm(uid, name, value string) error
	// Confirm submit async task
	Confirm(uid string) error
	// Cancel delete async task
	Cancel(uid string) error
	// Tasks get tasks
	Tasks(limit int) ([]*dto.TccAsyncTaskSchema, error)
	// Task get task
	Task(uid string) (*dto.TccAsyncTaskSchema, error)
}

// NewMysql ...
func NewMysql(db *gorm.DB) *mysql.TccAsyncTask {
	t := &mysql.TccAsyncTask{
		DB: db,
	}
	return t
}
