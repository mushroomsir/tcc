package lock

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/tcc/lock/mysql"
)

// LockInterface ...
type LockInterface interface {
	// Lock ...
	Lock(key string, expire time.Duration) error
	// Unlock ...
	Unlock(key string)
}

// NewMysql ...
func NewMysql(db *gorm.DB) *mysql.Lock {
	t := &mysql.Lock{
		DB: db,
	}
	return t
}
