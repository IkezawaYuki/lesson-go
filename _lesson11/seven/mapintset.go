package main

import (
	"bytes"
	"fmt"
	"sort"
)

type MapIntSet struct {
	set map[int]bool
}

func (s *MapIntSet) Has(x int) bool {
	if s.set == nil {
		return false
	}
	return s.set[x]
}

func (s *MapIntSet) Add(x int) {
	if s.set == nil {
		s.set = make(map[int]bool)
	}
	s.set[x] = true
}

func (s *MapIntSet) UnionWith(t *MapIntSet) {
	if t.set == nil {
		return
	}
	if s.set == nil {
		s.set = make(map[int]bool)
	}
	for x, b := range t.set {
		if b {
			s.set[x] = true
		}
	}
}

func (s *MapIntSet) String() string {
	if s.set == nil {
		return "{ }"
	}
	ints := make([]int, 0, len(s.set))
	for x, v := range s.set {
		if v {
			ints = append(ints, x)
		}
	}
	sort.Ints(ints)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, x := range ints {
		if i != 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", x)
	}
	buf.WriteByte('}')
	return buf.String()
}
