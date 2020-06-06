package mysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/logger/alog"
	"github.com/mushroomsir/tcc/lock/dto"
)

// Lock ...
type Lock struct {
	DB *gorm.DB
}

// Lock ...
func (a *Lock) Lock(key string, expire time.Duration) error {
	lock := &dto.Lock{}

	now := time.Now().UTC().Add(expire)
	insertSql := "insert into `tcc_lock` (`key`,`expire_at`) values (?, ?)"
	err := a.DB.Exec(insertSql, key, now).Error
	if err != nil {
		findSql := "select * from tcc_lock where `key` = ? limit 1"
		e := a.DB.Raw(findSql, key).Scan(lock).Error
		if e == nil {
			if lock.ExpireAt.Before(now) {
				a.Unlock(key) // release lock
				err = a.DB.Exec(insertSql, key, now).Error
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("%s locked, should expire at: %v, error: %s", key, expire, err.Error())
	}
	return err
}

// Unlock ...
func (a *Lock) Unlock(key string) {
	sql := "delete from tcc_lock where `key` = ?"
	err := a.DB.Exec(sql, key).Error
	if err != nil {
		alog.Errf("unlock: key %s, error %v", key, err.Error())
	}
}
