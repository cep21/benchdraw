package internal

// OrderedStringSet is a set of strings that remembers the order they are inserted in
type OrderedStringSet struct {
	Items map[string]struct{}
	Order []string
}

// Contains returns true if the set contains a string
func (s *OrderedStringSet) Contains(k string) bool {
	_, exists := s.Items[k]
	return exists
}

// Add a value to the set and append the key to Order.  Does nothing if the set already contains the key.
func (s *OrderedStringSet) Add(k string) {
	if s.Contains(k) {
		return
	}
	if s.Items == nil {
		s.Items = make(map[string]struct{})
	}
	s.Items[k] = struct{}{}
	s.Order = append(s.Order, k)
}
