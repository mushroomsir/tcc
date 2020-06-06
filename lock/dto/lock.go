package dto

import "time"

// Lock ...
type Lock struct {
	ID       int64     `gorm:"column:id"`
	Key      string    `gorm:"column:key"`
	ExpireAt time.Time `gorm:"column:expire_at"`
}
