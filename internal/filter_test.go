package internal

import (
	"testing"

	"github.com/cep21/benchparse"

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

func sameLists(t *testing.T, expected []benchparse.BenchmarkResult, seen []benchparse.BenchmarkResult) {
	t.Helper()
	require.Len(t, seen, len(expected))
	for i := 0; i < len(expected); i++ {
		require.Equal(t, expected[i].String(), seen[i].String())
		require.Equal(t, expected[i].Configuration, seen[i].Configuration)
	}
}

func TestFilter_FilterBenchmarks(t *testing.T) {
	filterList := func(data string, filter string, unit string, expected BenchmarkList) func(t *testing.T) {
		return func(t *testing.T) {
			bl := BenchmarkList(mustParse(data).Results)
			f := Filter{}
			got := f.FilterBenchmarks(bl, ToFilterPairs(filter), unit)
			sameLists(t, expected, got)
		}
	}
	t.Run("empty", filterList("", "", "", []benchparse.BenchmarkResult{}))
	t.Run("filterallgone", filterList(run1, "gone", "", []benchparse.BenchmarkResult{}))
	t.Run("filterempty", filterList(run1, "", "", []benchparse.BenchmarkResult{}))
	t.Run("simplefilter", filterList(run1, "BenchmarkTdigest_TotalSize", "ns/op", mustParse(`
goos: linux
goarch: amd64
pkg: github.com/cep21/tdigestbench
BenchmarkTdigest_TotalSize/digest=caio-8         	     100	  10687130 ns/op	   17920 B/op	      11 allocs/op
BenchmarkTdigest_TotalSize/digest=segmentio-8    	     422	   2844918 ns/op	 1380108 B/op	   54777 allocs/op
`).Results))
	t.Run("simplefilter", filterList(run1, "BenchmarkTdigest_TotalSize/digest=caio", "ns/op", mustParse(`
goos: linux
goarch: amd64
pkg: github.com/cep21/tdigestbench
BenchmarkTdigest_TotalSize/digest=caio-8         	     100	  10687130 ns/op	   17920 B/op	      11 allocs/op
`).Results))
	t.Run("multitestfilter", filterList(run1, "digest=caio", "ns/op", mustParse(`
goos: linux
goarch: amd64
pkg: github.com/cep21/tdigestbench
BenchmarkTdigest_TotalSize/digest=caio-8         	     100	  10687130 ns/op	   17920 B/op	      11 allocs/op
BenchmarkTdigest_Add/source=linear/digest=caio-8 	 1299776	       941 ns/op	      33 B/op	       0 allocs/op
BenchmarkTdigest_Add/source=rand/digest=caio-8                	 4080662	       317 ns/op	       0 B/op	       0 allocs/op
`).Results))
	t.Run("unitonone", filterList(run4, "BenchmarkTest", "allocs/op", mustParse(`
BenchmarkTest/name=bob/type=digest 1 10 ns/op 5 allocs/op
`).Results))
}
