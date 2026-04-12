package main

import (
	"bytes"
	"io"
	"testing"
)

type dummyReader struct {
	total   int64
	read    int64
	pattern byte
}

func (r *dummyReader) Read(p []byte) (n int, err error) {
	if r.read >= r.total {
		return 0, io.EOF
	}
	rem := r.total - r.read
	n = len(p)
	if int64(n) > rem {
		n = int(rem)
	}
	for i := range p[:n] {
		p[i] = r.pattern
	}
	r.read += int64(n)
	return n, nil
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size := 10 * 1024 * 1024
		payload := make([]byte, size)
		for j := range payload {
			payload[j] = 'B'
		}
		_ = bytes.NewReader(payload)
	}
}

func BenchmarkDummyReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size := int64(10 * 1024 * 1024)
		_ = &dummyReader{total: size, pattern: 'B'}
	}
}
