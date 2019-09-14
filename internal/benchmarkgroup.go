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

// ValuesByX returns all values in this group for a given x, ordered by all possible x values.  For example,
// if xDim='name', then allValues will contain all the values for xDim='name' and in the order we want to render
// them.  So if there are two names, Jack and John, then allValues=[Jack,John].  Unit is the benchmark value's unit
// that we search for and return (unit is the Y dimension drawn). If a benchmark group has no values for, then an empty
// array is returned for that index.  In other words, len(return_value) == len(alLValues)
func (b *BenchmarkGroup) ValuesByX(xDim string, unit string, allValues OrderedStringSet) [][]float64 {
	ret := make([][]float64, 0, len(allValues.Order))
	for _, v := range allValues.Order {
		allVals := make([]float64, 0, len(b.Results))
		for _, b := range b.Results {
			benchmarkKeys := makeKeys(b)
			if benchmarkKeys.Values[xDim] != v {
				continue
			}
			if val, exists := b.ValueByUnit(unit); exists {
				allVals = append(allVals, val)
			}
		}
		ret = append(ret, allVals)
	}
	return ret
}
