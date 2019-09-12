package internal

import (
	"github.com/cep21/benchparse"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func mustParse(s string) *benchparse.Run {
	ret, err := benchparse.Decoder{}.Decode(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return ret
}

const run1 = `
go test -v -benchmem -run=^$ -bench=. ./...
goos: linux
goarch: amd64
pkg: github.com/cep21/tdigestbench
BenchmarkTdigest_TotalSize/digest=caio-8         	     100	  10687130 ns/op	   17920 B/op	      11 allocs/op
BenchmarkTdigest_TotalSize/digest=segmentio-8    	     422	   2844918 ns/op	 1380108 B/op	   54777 allocs/op
BenchmarkTdigest_Add/source=linear/digest=caio-8 	 1299776	       941 ns/op	      33 B/op	       0 allocs/op
BenchmarkTdigest_Add/source=linear/digest=segmentio-8         	 1000000	      5602 ns/op	       8 B/op	       1 allocs/op
BenchmarkTdigest_Add/source=rand/digest=caio-8                	 4080662	       317 ns/op	       0 B/op	       0 allocs/op
BenchmarkTdigest_Add/source=rand/digest=segmentio-8           	 2220681	       785 ns/op	       8 B/op	       1 allocs/op
`


const run2 = `
go test -v -benchmem -run=^$ -bench=. ./...
name: john
unused: unused
BenchmarkTest/name=bob/type=digest 1 10 ns/op
BenchmarkTest/name=john 1 20 ns/op
type: sign
BenchmarkTest 1 30 ns/op
`


func makeMap(vals ...string) OrderedStringStringMap {
	var ret OrderedStringStringMap
	for i :=0;i<len(vals);i+=2 {
		ret.Insert(vals[i], vals[i+1])
	}
	return ret
}

func makeSet(vals ...string) OrderedStringSet {
	var ret OrderedStringSet
	for _, v := range vals {
		ret.Add(v)
	}
	return ret
}

func TestBenchmarkList_UniqueValuesForKey(t *testing.T) {
	bl := BenchmarkList(mustParse(run1).Results)
	require.Equal(t, makeSet("caio","segmentio"), bl.UniqueValuesForKey("digest"))
	require.Equal(t, makeSet(), bl.UniqueValuesForKey("bob"))
	b2 := BenchmarkList(mustParse(run2).Results)
	require.Equal(t, makeSet("bob", "john"), b2.UniqueValuesForKey("name"))
	require.Equal(t, makeSet("unused"), b2.UniqueValuesForKey("unused"))
	require.Equal(t, makeSet("unused"), b2.UniqueValuesForKey("unused"))
	require.Equal(t, makeSet("digest", "sign"), b2.UniqueValuesForKey("type"))
}