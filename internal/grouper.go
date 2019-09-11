package internal

type Grouper struct {
}

// each returned BenchmarkGroup will aggregate Results by unique groups Key/Value pairs
func (g *Grouper) GroupBenchmarks(in BenchmarkList, groups StringSet, unit string) BenchmarkGroupList {
	ret := make(BenchmarkGroupList, 0, len(in))
	setMap := make(map[string]*BenchmarkGroup)
	for _, b := range in {
		keysMap := MakeKeys(b)
		var hm HashableMap
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
