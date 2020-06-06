package tcc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToString(t *testing.T) {
	require := require.New(t)

	value := ObjToJSON("1")
	require.Equal("1", value)

	type test struct {
		key string
	}

	value = ObjToJSON(test{})
	require.Equal("{}", value)

	value = ObjToJSON(&test{})
	require.Equal("{}", value)
}
