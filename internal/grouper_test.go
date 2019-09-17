package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrouper_GroupBenchmarks(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		g := Grouper{}
		got := g.GroupBenchmarks(nil, makeSet())
		require.Equal(t, BenchmarkGroupList{}, got)
	})
	t.Run("run1", func(t *testing.T) {
		g := Grouper{}
		res := mustParse(run1).Results
		t.Run("simple", func(t *testing.T) {
			got := g.GroupBenchmarks(res, makeSet(""))
			require.Len(t, got, 1)
		})
		t.Run("firsttest", func(t *testing.T) {
			got := g.GroupBenchmarks(res, makeSet("BenchmarkTdigest_TotalSize"))
			require.Len(t, got, 2)
			require.Equal(t, "map[BenchmarkTdigest_TotalSize:]", got[0].Values.String())
			require.Equal(t, "map[]", got[1].Values.String())
		})
		t.Run("digest", func(t *testing.T) {
			got := g.GroupBenchmarks(res, makeSet("digest"))
			require.Len(t, got, 2)
			require.Equal(t, "map[digest:caio]", got[0].Values.String())
			require.Equal(t, "BenchmarkTdigest_TotalSize/digest=caio-8", got[0].Results[0].Name)
			require.Equal(t, "BenchmarkTdigest_Add/source=linear/digest=caio-8", got[0].Results[1].Name)
			require.Equal(t, "BenchmarkTdigest_Add/source=rand/digest=caio-8", got[0].Results[2].Name)
			require.Equal(t, "map[digest:segmentio]", got[1].Values.String())
		})
		t.Run("digest", func(t *testing.T) {
			got := g.GroupBenchmarks(res, makeSet("digest", "source"))
			require.Len(t, got, 6)
			require.Equal(t, "map[digest:caio]", got[0].Values.String())
			require.Equal(t, "map[digest:segmentio source:rand]", got[5].Values.String())
		})
	})
}
