package internal

import "github.com/cep21/benchparse"

type Filter struct {
}

type FilterPair struct {
	Key   string
	Value string
}

func (f *Filter) FilterBenchmarks(in []benchparse.BenchmarkResult, filters []FilterPair, unit string) []benchparse.BenchmarkResult {
	ret := make([]benchparse.BenchmarkResult, 0, len(in))
	for _, b := range in {
		// Benchmark must have a valid unit
		if _, exists := b.ValueByUnit(unit); !exists {
			continue
		}
		keys := b.AllKeyValuePairs()
		okToAdd := true
		// each filter must pass
		for _, f := range filters {
			val, exists := keys.Contents[f.Key]
			if !exists || (f.Value != "" && f.Value != val) {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			ret = append(ret, b)
		}
	}
	return ret
}