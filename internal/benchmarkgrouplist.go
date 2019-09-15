package internal

import "strings"

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
		if b[i].Values.Order[0] != expectedKey {
			return false
		}
	}
	return true
}

func (b BenchmarkGroupList) String() string {
	ret := make([]string, 0, len(b))
	for _, i := range b {
		ret = append(ret, i.String())
	}
	return "[" + strings.Join(ret, ",") + "]"
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
