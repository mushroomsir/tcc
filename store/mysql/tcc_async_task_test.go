package mysql

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // go-sql-driver
	"github.com/stretchr/testify/require"
)

func TestMysql(t *testing.T) {
	m := TccAsyncTask{DB: testDB}

	t.Run("try", func(t *testing.T) {
		require := require.New(t)

		uid := uuid.New().String()
		name := uuid.New().String()
		value := uuid.New().String()

		err := m.Try(uid, name, value)
		require.Nil(err)

		items, err := m.Tasks(1)
		require.Equal(uid, items[0].UID)
		require.Equal(1, len(items))
		require.Equal(name, items[0].Name)
		require.Equal(value, items[0].Value)
		require.Equal("try", items[0].Status)
		require.True(items[0].UpdatedAt.Unix() > 0)
		require.Equal(items[0].UpdatedAt, items[0].CreatedAt)

		err = m.Confirm(uid)
		require.Nil(err)

		items, err = m.Tasks(1)
		require.Equal("confirm", items[0].Status)

		m.Cancel(uid)

		items, err = m.Tasks(1)
		require.Equal(0, len(items))
	})

	t.Run("tryAndConfirm", func(t *testing.T) {
		require := require.New(t)

		uid := uuid.New().String()
		name := uuid.New().String()
		value := uuid.New().String()

		err := m.TryAndConfirm(uid, name, value)
		require.Nil(err)

		items, err := m.Tasks(1)
		require.Equal(uid, items[0].UID)
		require.Equal("confirm", items[0].Status)

		m.Cancel(uid)

		items, err = m.Tasks(1)
		require.Equal(0, len(items))
	})

	t.Run("task", func(t *testing.T) {
		require := require.New(t)
		_, err := m.Task(uuid.New().String())
		require.NotNil(err)

		uid := uuid.New().String()
		name := uuid.New().String()
		value := uuid.New().String()

		err = m.TryAndConfirm(uid, name, value)
		res, err := m.Task(uid)
		require.Nil(err)
		require.Equal(uid, res.UID)
		require.Equal("confirm", res.Status)
	})
}

// NewDB ...
func NewDB() *gorm.DB {
	query := "loc=UTC&readTimeout=10s&writeTimeout=10s&timeout=10s&multiStatements=true"
	parameters, err := url.ParseQuery(query)
	if err != nil {
		panic(err.Error())
	}

	parameters.Set("collation", "utf8mb4_general_ci")
	parameters.Set("parseTime", "true")

	user := "root"
	password := "password"
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
