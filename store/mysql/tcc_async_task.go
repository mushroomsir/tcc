package mysql

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/tcc/store/dto"
)

// TccAsyncTask ...
type TccAsyncTask struct {
	DB *gorm.DB
}

// TryPlan ...
func (a *TccAsyncTask) TryPlan(uid, name, value string) string {
	sql := a.genSql(uid, name, value, "try")
	return sql
}

// Try ....
func (a *TccAsyncTask) Try(uid, name, value string) error {
	sql := a.TryPlan(uid, name, value)
	return a.DB.Exec(sql).Error
}

// TryAndConfirm ....
func (a *TccAsyncTask) TryAndConfirm(uid, name, value string) error {
	sql := a.genSql(uid, name, value, "confirm")
	return a.DB.Exec(sql).Error
}

func (a *TccAsyncTask) genSql(uid, name, value, status string) string {
	sql := "insert into `tcc_async_task` (`uid`,`name`,`value`,`status`) values ('?','?','?','?')"
	sql = strings.Replace(sql, "?", uid, 1)
	sql = strings.Replace(sql, "?", name, 1)
	sql = strings.Replace(sql, "?", value, 1)
	sql = strings.Replace(sql, "?", status, 1)
	return sql
}

// Confirm ...
func (a *TccAsyncTask) Confirm(uid string) error {
	sql := "update tcc_async_task set status=?,updated_at=? where uid=?"
	return a.DB.Exec(sql, "confirm", time.Now().UTC(), uid).Error
}

// Cancel ...
func (a *TccAsyncTask) Cancel(uid string) error {
	sql := "delete from tcc_async_task where uid=?"
	return a.DB.Exec(sql, uid).Error
}

// Tasks ...
func (a *TccAsyncTask) Tasks(limit int) ([]*dto.TccAsyncTaskSchema, error) {
	sql := "select * from tcc_async_task order by id limit ?"
	tasks := []*dto.TccAsyncTaskSchema{}
	err := a.DB.Raw(sql, limit).Scan(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Task get task
func (a *TccAsyncTask) Task(uid string) (*dto.TccAsyncTaskSchema, error) {
	sql := "select * from tcc_async_task where uid=?"
	task := &dto.TccAsyncTaskSchema{}
	err := a.DB.Raw(sql, uid).Scan(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}
