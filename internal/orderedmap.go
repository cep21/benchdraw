package internal

import (
	"fmt"
	"sort"
	"strings"
)

// OrderedStringStringMap is a map that maintins insert order
type OrderedStringStringMap struct {
	Values map[string]string
	Order  []string
}

func mustWrite(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func mustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

// String pretty prints Values
func (h *OrderedStringStringMap) String() string {
	return fmt.Sprintf("%v", h.Values)
}

// Contains returns true if this value == this key and this key exists.  Useful if v can be empty
func (h *OrderedStringStringMap) Contains(k string, v string) bool {
	current, exists := h.Values[k]
	return exists && current == v
}

// Insert a value into the map
func (h *OrderedStringStringMap) Insert(k string, v string) {
	if _, exists := h.Values[k]; exists {
		h.Remove(k)
	}
	if h.Values == nil {
		h.Values = make(map[string]string)
	}
	h.Values[k] = v
	h.Order = append(h.Order, k)
}

// Remove a value from the map
func (h *OrderedStringStringMap) Remove(k string) {
	if h.Values == nil {
		return
	}
	delete(h.Values, k)
	for i, o := range h.Order {
		if o == k {
			h.Order = append(h.Order[:i], h.Order[i+1:]...)
			return
		}
	}
}

// Hash the contents of Values, ignoring order
func (h *OrderedStringStringMap) Hash() string {
	type kv struct {
		k string
		v string
	}
	toSort := make([]kv, 0, len(h.Values))
	for k, v := range h.Values {
		toSort = append(toSort, kv{k: k, v: v})
	}
	sort.Slice(toSort, func(i, j int) bool {
		return toSort[i].k < toSort[j].k
	})
	var uid strings.Builder
	for _, s := range toSort {
		if uid.Len() != 0 {
			mustWrite(uid.WriteString(string([]byte{0, 0})))
		}
		mustWrite(uid.WriteString(s.k))
		mustNotError(uid.WriteByte(0))
		mustWrite(uid.WriteString(s.v))
	}
	return uid.String()
}
