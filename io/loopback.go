package io

import (
	"io"

	"go.uber.org/multierr"
)

type LoopBack struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func (lb *LoopBack) Read(p []byte) (n int, err error) {
	return lb.r.Read(p)
}

func (lb *LoopBack) Write(p []byte) (n int, err error) {
	return lb.w.Write(p)
}

func (lb *LoopBack) Close() error {
	var err error
	if r := lb.r.Close(); err != nil {
		err = multierr.Append(err, r)
	}
	if w := lb.w.Close(); err != nil {
		err = multierr.Append(err, w)
	}
	return err
}

func NewLoopBack() io.ReadWriteCloser {
	lb := &LoopBack{}
	lb.r, lb.w = io.Pipe()
	return lb
}
