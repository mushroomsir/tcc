package dto

import "time"

// TccAsyncTaskSchema table `tcc_async_task`
type TccAsyncTaskSchema struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	UID    string `gorm:"column:uid"`
	Status string `gorm:"column:status"`
	Name   string `gorm:"column:name"`
	Value  string `gorm:"column:value"`
}
