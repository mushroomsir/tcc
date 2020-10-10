# TCC

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mushroomsir/tcc/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/mushroomsir/tcc)

TCC 基于本地消息表的分布式事务处理

## Design
1. [介绍分布式系统中的事务问题](https://github.com/mushroomsir/blog/blob/master/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F%E4%B8%AD%E7%9A%84%E4%BA%8B%E5%8A%A1%E9%97%AE%E9%A2%98.md)
2. TCC 异步补偿的设计与实现（TODO）

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



