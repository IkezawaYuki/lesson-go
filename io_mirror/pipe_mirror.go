package io_mirror

import (
	"errors"
	"sync"
)

type onceError struct {
	sync.Mutex
	err error
}

func (a *onceError) Store(err error) {
	a.Lock()
	defer a.Unlock()
	if a.err != nil {
		return
	}
	a.err = err
}

func (a *onceError) Load() error {
	a.Lock()
	defer a.Unlock()
	return a.err
}

var ErrClosedPipe = errors.New("io: read/write on closed pipe")

type pipe struct {
	wrMu sync.Mutex
	wrCh chan []byte
	rdCh chan int

	once sync.Once
	done chan struct{}
	rerr onceError
	werr onceError
}

func (p *pipe) Read(b []byte) (n int, err error) {
	select {
	case <-p.done:
		return 0, p.readCloseError()
	default:
	}

	select {
	case bw := <-p.wrCh:
		nr := copy(b, bw)
		p.rdCh <- nr
		return nr, nil
	case <-p.done:
		return 0, p.readCloseError()
	}
}

func (p *pipe) readCloseError() error {
	rerr := p.rerr.Load()
	if werr := p.werr.Load(); rerr == nil && werr != nil {
		return werr
	}
	return ErrClosedPipe
}

func (p *pipe) CloseRead(err error) error {
	if err == nil {
		err = ErrClosedPipe
	}
	p.rerr.Store(err)
	p.once.Do(func() {
		close(p.done)
	})
	return nil
}

func (p *pipe) Write(b []byte) (n int, err error) {
	select {
	case <-p.done:
		return 0, p.writeCloseError()
	default:
		p.wrMu.Lock()
		defer p.wrMu.Unlock()
	}

	for once := true; once || len(b) > 0; once = false {
		select {
		case p.wrCh <- b:
			nw := <-p.rdCh
			b = b[nw:]
			n += nw
		case <-p.done:
			return n, p.writeCloseError()
		}
	}
	return n, nil
}

func (p *pipe) writeCloseError() error {
	werr := p.werr.Load()
	if rerr := p.rerr.Load(); werr == nil && rerr != nil {
		return rerr
	}
	return ErrClosedPipe
}

func (p *pipe) CloseWrite(err error) error {
	if err == nil {
		err = EOF
	}
	p.werr.Store(err)
	p.once.Do(func() {
		close(p.done)
	})
	return nil
}

type PipeReader struct {
	p *pipe
}

func (r *PipeReader) Read(data []byte) (n int, err error) {
	return r.p.Read(data)
}

func (r *PipeReader) Close() error {
	return r.CloseWithError(nil)
}

func (r *PipeReader) CloseWithError(err error) error {
	return r.p.CloseRead(err)
}

type PipeWriter struct {
	p *pipe
}

func (w *PipeWriter) Write(data []byte) (n int, err error) {
	return w.p.Write(data)
}

func (w *PipeWriter) Close() error {
	return w.CloseWithError(nil)
}

func (w *PipeWriter) CloseWithError(err error) error {
	return w.p.CloseWrite(err)
}

func Pipe() (*PipeReader, *PipeWriter) {
	p := &pipe{
		wrCh: make(chan []byte),
		rdCh: make(chan int),
		done: make(chan struct{}),
	}
	return &PipeReader{p}, &PipeWriter{p}
}
