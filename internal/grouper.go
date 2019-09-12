package internal

// Grouper understands how to combine benchmarks into groups of common key/value pairs
type Grouper struct {
}

// GroupBenchmarks will return for each returned BenchmarkGroup will aggregate Results by unique groups Key/Value pairs
func (g *Grouper) GroupBenchmarks(in BenchmarkList, groups OrderedStringSet, unit string) BenchmarkGroupList {
	ret := make(BenchmarkGroupList, 0, len(in))
	setMap := make(map[string]*BenchmarkGroup)
	for _, b := range in {
		keysMap := makeKeys(b)
		var hm OrderedStringStringMap
		if len(groups.Order) == 0 {
			// Group by everything except unit
			for _, k := range keysMap.Order {
				if k != unit {
					hm.Insert(k, keysMap.Values[k])
				}
			}
		} else {
			for _, ck := range groups.Order {
				if configValue, exists := keysMap.Values[ck]; exists {
					hm.Insert(ck, configValue)
				}
			}
		}
		mapHash := hm.Hash()
		if existing, exists := setMap[mapHash]; exists {
			existing.Results = append(existing.Results, b)
		} else {
			bg := &BenchmarkGroup{
				Values:  hm,
				Results: BenchmarkList{b},
			}
			setMap[mapHash] = bg
			ret = append(ret, bg)
		}
	}
	return ret
}
