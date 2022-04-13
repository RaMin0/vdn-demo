package main

import (
	"fmt"
	stdlog "log"
	"sync"
	"time"
)

func init() {
	stdlog.SetFlags(0)
}

var log = &tickerLogger{}

type tickerLogger struct {
	s  time.Time
	p  int
	mx sync.Mutex
}

func (l *tickerLogger) Printf(format string, args ...interface{}) {
	if l.s.IsZero() {
		l.s = time.Now()
	}

	l.mx.Lock()
	s := int(time.Since(l.s).Seconds())
	if s > l.p {
		stdlog.Print()
		l.p = s
	}
	l.mx.Unlock()

	stdlog.Printf("[%3d] %s", s, fmt.Sprintf(format, args...))
}
