package internal

type StringSet struct {
	Items map[string]struct{}
	Order []string
}

func (s *StringSet) Contains(k string) bool {
	_, exists := s.Items[k]
	return exists
}

func (s *StringSet) Add(k string) {
	if s.Contains(k) {
		return
	}
	if s.Items == nil {
		s.Items = make(map[string]struct{})
	}
	s.Items[k] = struct{}{}
	s.Order = append(s.Order, k)
}

