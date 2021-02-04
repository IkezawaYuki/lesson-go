package hash_error

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
