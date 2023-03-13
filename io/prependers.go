package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
)

type Prepender struct {
	color  string
	prefix string
	w      io.Writer
}

var ErrInvalidColor = errors.New("invalid color, or not supported")

func NewPrepender(w io.Writer, prefix string, color string) (*Prepender, error) {
	if _, ok := tcell.ColorNames[color]; !ok {
		return nil, ErrInvalidColor
	}

	return &Prepender{color: color, prefix: prefix, w: w}, nil
}

func (p *Prepender) Write(b []byte) (int, error) {
	return p.w.Write(append([]byte(fmt.Sprintf("[%s]%s[white] ", p.color, p.prefix)), append(bytes.TrimSpace(b), []byte("\n")...)...))
}

// NewRecvPrepender returns a new RecvPrepender wrapping the given writer.
func NewRecvPrepender(w io.Writer) *Prepender {
	p, _ := NewPrepender(w, "RECV:", "orange")
	return p
}

// NewSendPrepender returns a new Prepender, and is meant to be used in conjunction with a tview.TextView. It prepends
// any text written with a
// tcell.Color formatted string "[yellow]SEND:[white] ". It is particularly useful when used sibling to a serial device
// writer as a parameter to an io.MultiWriter to echo a sent command in the text view.
// NewSendPrepender returns a new SendPrepender wrapping the given writer.
func NewSendPrepender(w io.Writer) *Prepender {
	p, _ := NewPrepender(w, "SEND:", "yellow")
	return p
}
