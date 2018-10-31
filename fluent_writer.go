package main

import (
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type FluentWriter struct {
	f   *fluent.Fluent
	tag string
}

func NewFluentWriter(f *fluent.Fluent, tag string) *FluentWriter {
	return &FluentWriter{
		f:   f,
		tag: tag,
	}
}

func (w *FluentWriter) Post(ts time.Time, sample sample) error {
	return w.f.PostWithTime(w.tag, ts, sample)
}
