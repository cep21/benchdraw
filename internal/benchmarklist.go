package internal

import "github.com/cep21/benchparse"

// BenchmarkList is a list of benchmarks
type BenchmarkList []benchparse.BenchmarkResult

// UniqueValuesForKey returns all common value properties of each benchmark in this list
// for a key.
func (b BenchmarkList) UniqueValuesForKey(key string) OrderedStringSet {
	var ret OrderedStringSet
	for _, b := range b {
		keys := makeKeys(b)
		if keyValue, exists := keys.Values[key]; exists {
			ret.Add(keyValue)
		}
	}
	return ret
}

func makeKeys(r benchparse.BenchmarkResult) OrderedStringStringMap {
	nameKeys := r.AllKeyValuePairs()
	var ret OrderedStringStringMap
	for _, k := range nameKeys.Order {
		ret.Insert(k, nameKeys.Contents[k])
	}
	return ret
}

// ValuesByX returns all values in this group for a given x, ordered by all possible x values.  For example,
// if xDim='name', then allValues will contain all the values for xDim='name' and in the order we want to render
// them.  So if there are two names, Jack and John, then allValues=[Jack,John].  Unit is the benchmark value's unit
// that we search for and return (unit is the Y dimension drawn). If a benchmark group has no values for, then an empty
// array is returned for that index.  In other words, len(return_value) == len(alLValues)
func (b BenchmarkList) ValuesByX(xDim string, unit string, allValues OrderedStringSet) [][]float64 {
	ret := make([][]float64, 0, len(allValues.Order))
	for _, v := range allValues.Order {
		allVals := make([]float64, 0, len(b))
		for _, b := range b {
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
