package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBenchmarkGroupList_AllSingleKey(t *testing.T) {
	require.True(t, BenchmarkGroupList([]*BenchmarkGroup{
		{
			Values: makeMap("name", "bob"),
		},
		{
			Values: makeMap("name", "john"),
		},
	}).AllSingleKey())
	require.True(t, BenchmarkGroupList([]*BenchmarkGroup{}).AllSingleKey())
	require.False(t, BenchmarkGroupList([]*BenchmarkGroup{
		{
			Values: makeMap("name", "bob", "age", "12"),
		},
		{
			Values: makeMap("name", "bob"),
		},
	}).AllSingleKey())
	require.False(t, BenchmarkGroupList([]*BenchmarkGroup{
		{
			Values: makeMap("age", "12"),
		},
		{
			Values: makeMap("name", "bob"),
		},
	}).AllSingleKey())
}

func TestBenchmarkGroupList_Normalize(t *testing.T) {
	testNormalEqual := func(bgl BenchmarkGroupList, expected BenchmarkGroupList) func(t *testing.T) {
		return func(t *testing.T) {
			bgl.Normalize()
			require.Equal(t, expected.String(), bgl.String())
		}
	}
	t.Run("empty", testNormalEqual(BenchmarkGroupList(nil), BenchmarkGroupList(nil)))
	t.Run("simple", testNormalEqual([]*BenchmarkGroup{
		{
			Values: makeMap("age", "12"),
		},
	}, []*BenchmarkGroup{
		{
			Values: makeMap(),
		},
	}))
	t.Run("basicremove", testNormalEqual([]*BenchmarkGroup{
		{
			Values: makeMap("age", "12", "name", "john"),
		},
		{
			Values: makeMap("age", "13", "name", "john"),
		},
	}, []*BenchmarkGroup{
		{
			Values: makeMap("age", "12"),
		},
		{
			Values: makeMap("age", "13"),
		},
	}))
	t.Run("same", testNormalEqual([]*BenchmarkGroup{
		{
			Values: makeMap("age", "12", "name", "john"),
		},
		{
			Values: makeMap("age", "13", "name", "jack"),
		},
	}, []*BenchmarkGroup{
		{
			Values: makeMap("age", "12", "name", "john"),
		},
		{
			Values: makeMap("age", "13", "name", "jack"),
		},
	}))
	t.Run("onlytwomatch", testNormalEqual([]*BenchmarkGroup{
		{
			Values: makeMap("name", "john"),
		},
		{
			Values: makeMap("name", "john"),
		},
		{
			Values: makeMap("name", "jack"),
		},
	}, []*BenchmarkGroup{
		{
			Values: makeMap("name", "john"),
		},
		{
			Values: makeMap("name", "john"),
		},
		{
			Values: makeMap("name", "jack"),
		},
	}))
}
