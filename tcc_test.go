package tcc

import (
	"sync"
	"testing"

	"github.com/mushroomsir/tcc/lock"
	"github.com/mushroomsir/tcc/store"
	"github.com/stretchr/testify/require"
)

func TestTcc(t *testing.T) {

	tc := New(&Option{
		PullTaskInterval: 1,
		Store:            store.NewMysql(testDB),
		Lock:             lock.NewMysql(testDB),
	})
	t.Run("try", func(t *testing.T) {
		require := require.New(t)

		name := "send.task.create"
		value := "value"

		tx := tc.NewTransaction(name)
		err := tx.Try(value)
		require.Nil(err)

		items, err := tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(name, items[0].Name)
		require.Equal(value, items[0].Value)
		require.Equal("try", items[0].Status)
		require.True(items[0].UpdatedAt.Unix() > 0)
		require.Equal(items[0].UpdatedAt, items[0].CreatedAt)

		err = tx.Cancel()
		require.Nil(err)

		items, err = tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(0, len(items))
	})

	t.Run("tryAndConfirm", func(t *testing.T) {
		require := require.New(t)

		name := "send.task.update"
		value := "value"

		tx := tc.NewTransaction(name)
		err := tx.TryAndConfirm(value)
		require.Nil(err)

		items, err := tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(name, items[0].Name)
		require.Equal(value, items[0].Value)
		require.Equal("confirm", items[0].Status)
		require.True(items[0].UpdatedAt.Unix() > 0)
		require.Equal(items[0].UpdatedAt, items[0].CreatedAt)

		tx.Cancel()

		items, err = tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(0, len(items))
	})

	t.Run("loop try", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		wg.Add(1)
		tc.SetTryHandler(func(t *Task) {
			wg.Done()
		})

		name := "send.task.update"
		value := "value"

		tx := tc.NewTransaction(name)
		err := tx.Try(value)
		require.Nil(err)

		items, err := tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(1, len(items))
		require.Equal(tccTry, items[0].Status)

		wg.Wait()
		tx.Cancel()

		items, err = tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(0, len(items))
	})

	t.Run("loop confirm", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		wg.Add(1)
		tc.SetConfirmHandler(func(t *Task) {
			err := t.Cancel()
			require.Nil(err)
			wg.Done()
		})

		name := "send.task.update"
		value := "value"

		tx := tc.NewTransaction(name)
		err := tx.TryAndConfirm(value)
		require.Nil(err)

		items, err := tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(1, len(items))
		require.Equal(tccConfirm, items[0].Status)

		wg.Wait()

		items, err = tc.Store.Tasks(1)
		require.Nil(err)
		require.Equal(0, len(items))
	})
}
