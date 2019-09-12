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
