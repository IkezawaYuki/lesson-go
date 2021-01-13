package intset

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var seed int64 = time.Now().UTC().UnixNano()

func BenchmarkIntSet_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add_IntSet()
	}
}

func add_IntSet() {
	var set IntSet
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 500; i++ {
		set.Add(rng.Intn(math.MaxInt16))
	}
}

func add_MapIntSet() {
	var set MapIntSet
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 500; i++ {
		set.Add(rng.Intn(math.MaxInt16))
	}
}
