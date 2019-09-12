package internal

import (
	"strings"

	"github.com/cep21/benchparse"
)

// Filter understands how to filter a benchmark result set
type Filter struct {
}

// FilterPair controls how to filter.  It means filter only key=value.  If value is empty, then filters for existence.
type FilterPair struct {
	Key   string
	Value string
}

// ToFilterPairs converts a string of benchmark format to a list of filter pairs.  For example, BenchmarkBob/name=john
// would become [{Key: BenchmarkBob}, {Key: name, Value: bob}]
func ToFilterPairs(s string) []FilterPair {
	parts := strings.Split(s, "/")
	ret := make([]FilterPair, 0, len(parts))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 1 {
			ret = append(ret, FilterPair{
				Key: p,
			})
		} else {
			ret = append(ret, FilterPair{
				Key:   kv[0],
				Value: kv[1],
			})
		}
	}
	return ret
}

// FilterBenchmarks only returns benchmarks that contain this unit and belong to the filter pairs.
func (f *Filter) FilterBenchmarks(in []benchparse.BenchmarkResult, filters []FilterPair, unit string) BenchmarkList {
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
