package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makePairs(args ...string) []FilterPair {
	ret := make([]FilterPair, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		ret = append(ret, FilterPair{
			Key:   args[i],
			Value: args[i+1],
		})
	}
	return ret
}

func TestToFilterPairs(t *testing.T) {
	filtersEqual := func(arg string, want []FilterPair) func(t *testing.T) {
		return func(t *testing.T) {
			require.Equal(t, want, ToFilterPairs(arg))
		}
	}
	t.Run("empty", filtersEqual("", makePairs()))
	t.Run("simple", filtersEqual("simple", makePairs("simple", "")))
	t.Run("set", filtersEqual("key=value", makePairs("key", "value")))
	t.Run("settwo", filtersEqual("key=value/key2=value2", makePairs("key", "value", "key2", "value2")))
	t.Run("setwithempty", filtersEqual("key=value/", makePairs("key", "value")))
}
