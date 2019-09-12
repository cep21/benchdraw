package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderedStringStringMap(t *testing.T) {
	var m OrderedStringStringMap
	require.False(t, m.Contains("", ""))
	m.Insert("name", "bob")
	require.False(t, m.Contains("", ""))
	require.False(t, m.Contains("name", "john"))
	require.True(t, m.Contains("name", "bob"))
	h1 := m.Hash()
	m.Insert("name", "john")
	h2 := m.Hash()
	require.True(t, m.Contains("name", "john"))
	require.False(t, m.Contains("name", "bob"))
	require.NotEqual(t, h1, h2)
	m.Remove("name")
	require.Equal(t, []string{}, m.Order)
}
