package internal

import (
	"fmt"
	"strings"
)

// BenchmarkGroup contains a set of parsed benchmark results for a group of values.  All benchmarks will have the
// key/value pairs represented in Values
type BenchmarkGroup struct {
	// Values are the common key/value pairs that make up this benchmark group
	Values OrderedStringStringMap
	// Results are a list of benchmark results for this group
	Results BenchmarkList
}

func (b *BenchmarkGroup) String() string {
	return fmt.Sprintf("vals=%v len_results=%d", b.Values, len(b.Results))
}

// NominalLineName returns how we should render this BenchmarkGroup's name in the UI.  If they are all a single key
// then we just return the first value.  Otherwise, return [key1=value1,key2=value2,...]
func NominalLineName(values OrderedStringStringMap, singleKey bool) string {
	if singleKey && len(values.Order) > 0 {
		return values.Values[values.Order[0]]
	}
	ret := make([]string, 0, len(values.Order))
	for _, c := range values.Order {
		ret = append(ret, c+"="+values.Values[c])
	}
	if len(ret) == 0 {
		return ""
	}
	return "[" + strings.Join(ret, ",") + "]"
}
