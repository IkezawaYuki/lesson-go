package main

import (
	"bytes"
	"fmt"
)

const UINT_SIZE = 32 << (^uint(0) >> 63)

type IntSet struct {
	words []uint
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/UINT_SIZE, uint(x%UINT_SIZE)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/UINT_SIZE, uint(x%UINT_SIZE)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) AddAll(values ...int) {
	for _, x := range values {
		s.Add(x)
	}
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		}
	}
	for i := len(t.words); i < len(s.words); i++ {
		s.words[i] = 0
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) Enums() []int {
	len := s.Len()
	if len == 0 {
		return []int{}
	}
	enums := make([]int, 0, len)
	for i, sword := range s.words {
		for bit := uint(0); bit < UINT_SIZE; bit++ {
			if sword&(1<<bit) != 0 {
				enums = append(enums, i*UINT_SIZE+int(bit))
			}
		}
	}
	return enums
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < UINT_SIZE; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", UINT_SIZE*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *IntSet) Len() int {
	len := 0
	for _, word := range s.words {
		len += bitCount(word)
	}
	return len
}

func bitCount(x uint) int {
	if UINT_SIZE == 32 {
		return bitCount32(uint32(x))
	} else {
		return bitCount64(uint64(x))
	}
}

func bitCount64(x uint64) int {
	// Hacker's Delight, Figure 5-2.
	x = x - ((x >> 1) & 0x5555555555555555)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0f0f0f0f0f0f0f0f
	x = x + (x >> 8)
	x = x + (x >> 16)
	x = x + (x >> 32)
	return int(x & 0x7f)
}

func bitCount32(x uint32) int {
	// Hacker's Delight, Figure 5-2.
	x = x - ((x >> 1) & 0x55555555)
	x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
	x = (x + (x >> 4)) & 0x0f0f0f0f
	x = x + (x >> 8)
	x = x + (x >> 16)
	return int(x & 0x7f)
}

func (s *IntSet) Remove(x int) {
	word, bit := x/UINT_SIZE, uint(x%UINT_SIZE)
	if word > len(s.words) {
		return
	}
	s.words[word] &^= 1 << bit
}

func (s *IntSet) Clear() {
	s.words = []uint{}
}

func (s *IntSet) Copy() *IntSet {
	var c IntSet

	c.words = make([]uint, len(s.words))
	copy(c.words, s.words)
	return &c
}
