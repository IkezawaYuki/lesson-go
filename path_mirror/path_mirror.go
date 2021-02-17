package path_mirror

import "strings"

type lazybuf struct {
	s   string
	buf []byte
	w   int
}

func (b *lazybuf) index(i int) byte {
	if b.buf != nil {
		return b.buf[i]
	}
	return b.s[i]
}

func (b *lazybuf) append(c byte) {
	if b.buf == nil {
		if b.w < len(b.s) && b.s[b.w] == c {
			b.w++
			return
		}
		b.buf = make([]byte, len(b.s))
		copy(b.buf, b.s[:b.w])
	}
	b.buf[b.w] = c
	b.w++
}

func (b *lazybuf) string() string {
	if b.buf == nil {
		return b.s[:b.w]
	}
	return string(b.buf[:b.w])
}

func Clean(path string) string {
	if path == "" {
		return "."
	}

	rooted := path[0] == '/'
	n := len(path)

	out := lazybuf{s: path}
	r, dotdot := 0, 0
	if rooted {
		out.append('/')
	}

	for r < n {
		switch {
		case path[r] == '/':
			r++
		case path[r] == '.' && (r+1 == n || path[r+1] == '/'):
			r++
		case path[r] == '.' && path[r+1] == '.' && (r+2 == n || path[r+2] == '/'):
			r += 2
			switch {
			case out.w > dotdot:
				out.w--
				for out.w > dotdot && out.index(out.w) != '/' {
					out.w--
				}
			case !rooted:
				if out.w > 0 {
					out.append('/')
				}
				out.append('.')
				out.append('.')
				dotdot = out.w
			}
		default:
			if rooted && out.w != 1 || !rooted && out.w != 0 {
				out.append('/')
			}
			for ; r < n && path[r] != '/'; r++ {
				out.append(path[r])
			}
		}
	}

	if out.w == 0 {
		return "."
	}
	return out.string()
}

func Split(path string) (dir, file string) {
	i := strings.LastIndex(path, "/")
	return path[:i+1], path[i+1:]
}

func Join(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return Clean(strings.Join(elem[i:], "/"))
		}
	}
	return ""
}

func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

func Base(path string) string {
	if path == "" {
		return "."
	}
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}

	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}

	if path == "" {
		return "/"
	}
	return path
}

func IsAbs(path string) bool {
	return len(path) > 0 && path[0] == '/'
}

func Dir(path string) string {
	dir, _ := Split(path)
	return Clean(dir)
}
