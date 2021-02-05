package hash_error

import "context"

type Seed struct {
	s uint64
}

type Hash struct {
	_     [0]func()
	seed  Seed
	state Seed
	buf   [64]byte
	n     int
}

func (h *Hash) initSeed() {
	if h.seed.s == 0 {
		h.setSeed(MakeSeed())
	}
}

func (h *Hash) Write(b []byte) (int, error) {
	size := len(b)
	for h.n+len(b) > len(h.buf) {
		k := copy(h.buf[h.n:], b)
		h.n = len(h.buf)
		b = b[k:]
		h.flush()
	}
	h.n += copy(h.buf[h.n:], b)
	return size, nil
}

func (h *Hash) WriteString(s string) (int, error) {
	size := len(s)
	for h.n+len(s) > len(h.buf) {
		k := copy(h.buf[h.n:], s)
		h.n = len(h.buf)
		s = s[k:]
		h.flush()
	}
	h.n += copy(h.buf[h.n:], s)
	return size, nil
}

func (h *Hash) Seed() Seed {
	h.initSeed()
	return h.seed
}

func (h *Hash) setSeed(seed Seed) {
	if seed.s == 0 {
		panic("maphash: use of uninitialized Seed")
	}
	h.seed = seed
	h.state = seed
}

func (h *Hash) flush() {
	if h.n != len(h.buf) {
		panic("maphash: flush of partially full buffer")
	}
	h.initSeed()
	h.state.s = rthash(h.buf[:], h.state.s)
	h.n = 0
}

func (h *Hash) Sum64() uint64 {
	h.initSeed()
	return rthash(h.buf[:h.n], h.state.s)
}

func MakeSeed() Seed {
	var s1, s2 uint64
	for {
		s1 = uint64(runtime_fastrand())
		s2 = uint64(runtime_fastrand())
		if s1|s2 != 0 {
			break
		}
	}
	return Seed{s: s1<<32 + s2}
}

func runtime_fastrand() uint32 {}

func (h *Hash) Sum(b []byte) []byte {
	x := h.Sum64()
	return append(b,
		byte(x>>0),
		byte(x>>8),
		byte(x>>16),
		byte(x>>32),
		byte(x>>40),
		byte(x>>48),
		byte(x>>56),
	)
	context.Background()
}

func (h *Hash) Size() int { return 8 }

func (h *Hash) BlockSize() int { return len(h.buf) }
