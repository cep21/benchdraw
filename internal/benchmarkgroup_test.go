package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNominalLineName(t *testing.T) {
	require.Equal(t, NominalLineName(makeMap("name", "john"), true), "john")
	require.Equal(t, NominalLineName(makeMap("name", "john"), false), "[name=john]")
	require.Equal(t, NominalLineName(makeMap("name", "john", "age", "12"), true), "john")
	require.Equal(t, NominalLineName(makeMap("name", "john", "age", "12"), false), "[name=john,age=12]")
}
