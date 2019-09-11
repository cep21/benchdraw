package internal

import "github.com/cep21/benchparse"

type BenchmarkList []benchparse.BenchmarkResult

func (b BenchmarkList) UniqueValuesForKey(key string) StringSet {
	var ret StringSet
	for _, b := range b {
		keys := makeKeys(b)
		if keyValue, exists := keys.Values[key]; exists {
			ret.Add(keyValue)
		}
	}
	return ret
}

func makeKeys(r benchparse.BenchmarkResult) HashableMap {
	nameKeys := r.AllKeyValuePairs()
	var ret HashableMap
	for _, k := range nameKeys.Order {
		ret.Insert(k, nameKeys.Contents[k])
	}
	return ret
}
