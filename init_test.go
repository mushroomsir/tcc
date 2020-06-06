package tcc

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // go-sql-driver
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	testDB = NewDB()
	testDB.Exec("delete from tcc_async_task")
	os.Exit(m.Run())
}

// NewDB ...
func NewDB() *gorm.DB {
	query := "loc=UTC&readTimeout=10s&writeTimeout=10s&timeout=10s&multiStatements=true"
	parameters, err := url.ParseQuery(query)
	if err != nil {
		panic(err.Error())
	}
	// 强制使用
	parameters.Set("collation", "utf8mb4_general_ci")
	parameters.Set("parseTime", "true")

	user := "root"
	password := "root"
	host := "localhost:3306"
	database := "test_tcc"
	// https://github.com/go-sql-driver/mysql#parameters
	url := fmt.Sprintf(`%s:%s@(%s)/%s?%s`, user, password, host, database, parameters.Encode())
	db, err := gorm.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}
	db.SingularTable(true)
	db.LogMode(false)
	db.DB().SetMaxIdleConns(8)
	db.DB().SetMaxOpenConns(64)
	return db
}
