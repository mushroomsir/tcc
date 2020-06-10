package mysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/tcc/lock/dto"
)

// Lock ...
type Lock struct {
	DB *gorm.DB
}

// Lock ...
func (a *Lock) Lock(name string, expire time.Duration) error {
	lock := &dto.Lock{}

	now := time.Now().UTC().Add(expire)
	insertSql := "insert into `tcc_lock` (`name`,`expire_at`) values (?, ?)"
	err := a.DB.Exec(insertSql, name, now).Error
	if err != nil {
		findSql := "select * from tcc_lock where `name` = ? limit 1"
		err = a.DB.Raw(findSql, name).Scan(lock).Error
		if err == nil {
			if lock.ExpireAt.Before(now) {
				a.Unlock(name) // release lock
				err = a.DB.Exec(insertSql, name, now).Error
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("%s locked, should expire at: %v, error: %s", name, expire, err.Error())
	}
	return err
}

// Unlock ...
func (a *Lock) Unlock(name string) error {
	sql := "delete from tcc_lock where `name` = ?"
	err := a.DB.Exec(sql, name).Error
	if err != nil {
		return fmt.Errorf("unlock: key %s, error %s", name, err.Error())
	}
	return nil
}
