# TCC

[![Build Status](https://img.shields.io/travis/mushroomsir/tcc.svg?style=flat-square)](https://travis-ci.org/mushroomsir/tcc)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/tcc.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/tcc?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mushroomsir/tcc/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/mushroomsir/tcc)

TCC 分布式系统中的异步补偿

## Usage
```sh
go get github.com/mushroomsir/tcc
```
Create transaction table by `sql/mysql.sql`.

### Demo

```go
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
```



