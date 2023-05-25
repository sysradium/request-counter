package decorators

import (
	"bufio"
	"context"
	"encoding/gob"
	"os"
	"time"

	"github.com/sysradium/request-counter/internal/counter/ephemeral"
)

const defaultBuffSize = 512

// Persisted decorator records each request to a file, buffering writes
type Persisted struct {
	*ephemeral.EphemeralSlidingStorage
	buf     *bufio.Writer
	encoder *gob.Encoder
}

func (p *Persisted) Add(t time.Time) error {
	if err := p.EphemeralSlidingStorage.Add(t); err != nil {
		return err
	}
	if err := p.encoder.Encode(t); err != nil {
		return err
	}
	return nil
}

func (p *Persisted) Start(ctx context.Context, done chan struct{}) error {
	go func() {
		t := time.NewTicker(time.Second * 1)
		defer t.Stop()

	loop:
		for {
			select {
			case <-t.C:
				p.Flush()
			case <-ctx.Done():
				p.Flush()
				break loop
			}
		}
		<-done

	}()
	return p.EphemeralSlidingStorage.Start()
}

func (p *Persisted) Flush() {
	p.buf.Flush()
}

func NewPersisted(
	s *ephemeral.EphemeralSlidingStorage,
	f *os.File,
) *Persisted {
	buffer := bufio.NewWriterSize(f, defaultBuffSize)

	return &Persisted{
		buf:                     buffer,
		encoder:                 gob.NewEncoder(buffer),
		EphemeralSlidingStorage: s,
	}
}
