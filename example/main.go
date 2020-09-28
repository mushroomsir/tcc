package main

import (
	"github.com/jinzhu/gorm"
	"github.com/mushroomsir/tcc"
	"github.com/mushroomsir/tcc/store"
)

func main() {
	tc := tcc.New(&tcc.Option{
		PullTaskInterval: 1,
		Store:            store.NewMysql(&gorm.DB{}),
	})

	tx := tc.NewTransaction("name")
	sql := tx.TryPlan("value")

	err := doSomeThing(sql)

	if err != nil {
		tx.Confirm() // confirm to summit async compensation task
	} else {
		tx.Cancel() // cancel async compensation task
	}
}

func doSomeThing(sql string) error {
	// execute sql in transaction
	return nil
}
