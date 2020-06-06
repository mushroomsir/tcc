package tcc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToString(t *testing.T) {
	require := require.New(t)

	value := ObjToJson("1")
	require.Equal("1", value)

	type test struct {
		key string
	}

	value = ObjToJson(test{})
	require.Equal("{}", value)

	value = ObjToJson(&test{})
	require.Equal("{}", value)
}
