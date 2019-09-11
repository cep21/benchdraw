package internal

import (
	"fmt"
	"strings"
)

type BenchmarkGroup struct {
	Values  HashableMap
	Results BenchmarkList
}

type BenchmarkGroupList []*BenchmarkGroup

func (b *BenchmarkGroup) String() string {
	return fmt.Sprintf("vals=%v len_results=%d", b.Values, len(b.Results))
}

func (b *BenchmarkGroup) NominalLineName(singleKey bool) string {
	if singleKey && len(b.Values.Order) > 0 {
		return b.Values.Values[b.Values.Order[0]]
	}
	ret := make([]string, 0, len(b.Values.Order))
	for _, c := range b.Values.Order {
		ret = append(ret, c+"="+b.Values.Values[c])
	}
	if len(ret) == 0 {
		return ""
	}
	return "[" + strings.Join(ret, ",") + "]"
}
