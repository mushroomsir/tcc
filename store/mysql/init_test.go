package mysql

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	testDB = NewDB()
	testDB.Exec("delete from tcc_async_task")
	os.Exit(m.Run())
}
