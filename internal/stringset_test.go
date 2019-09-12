package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderedStringSet(t *testing.T) {
	var s OrderedStringSet
	require.False(t, s.Contains(""))
	s.Add("bob")
	require.False(t, s.Contains("john"))
	require.True(t, s.Contains("bob"))
	s.Add("name")
	require.True(t, s.Contains("name"))
	require.Equal(t, []string{"bob", "name"}, s.Order)
	require.True(t, s.Contains("bob"))
	require.Equal(t, []string{"bob", "name"}, s.Order)
}
