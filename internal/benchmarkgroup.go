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

type BenchmarkGroupList []*BenchmarkGroup

// AllSingleKey returns true if all the benchmarks in this group are of a single value.  This can help us render
// [name=bob] as just bob in the UI
func (b BenchmarkGroupList) AllSingleKey() bool {
	if len(b) <= 1 {
		return true
	}
	if len(b[0].Values.Order) > 1 {
		return false
	}
	expectedKey := b[0].Values.Order[0]
	for i := 1; i < len(b); i++ {
		if len(b[i].Values.Order) > 1 {
			return false
		}
		if b[0].Values.Order[0] != expectedKey {
			return false
		}
	}
	return true
}

// Normalize removes all values that are the same in each BenchmarkGroup.  For example, if all benchmarks have the
// value os=darwin, then we remove os=darwin from the BenchmarkGroup's Values map.
func (b BenchmarkGroupList) Normalize() {
	if len(b) == 0 {
		return
	}
	keysToRemove := make([]string, 0, len(b[0].Values.Values))
	for k, v := range b[0].Values.Values {
		canRemoveValue := true
	checkRestLoop:
		for i := 1; i < len(b); i++ {
			if !b[i].Values.Contains(k, v) {
				canRemoveValue = false
				break checkRestLoop
			}
		}
		if canRemoveValue {
			keysToRemove = append(keysToRemove, k)
		}
	}
	for _, k := range keysToRemove {
		for _, i := range b {
			i.Values.Remove(k)
		}
	}
}

func (b *BenchmarkGroup) String() string {
	return fmt.Sprintf("vals=%v len_results=%d", b.Values, len(b.Results))
}

// NominalLineName returns how we should render this BenchmarkGroup's name in the UI.  If they are all a single key
// then we just return the first value.  Otherwise, return [key1=value1,key2=value2,...]
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
