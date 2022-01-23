package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStrPtr(t *testing.T) {
	value := "test"
	ptr := strptr(value)
	require.Equal(t, value, *ptr)
}
