package snapshot

import (
	"context"
	"time"
)

type getter interface{ Get() []time.Time }
type flusher func([]byte)
type marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

type Periodic struct {
	period  time.Duration
	store   getter
	flush   flusher
	encoder marshaler
}

func (p *Periodic) Run(ctx context.Context) {
	t := time.NewTimer(p.period)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			d := p.store.Get()
			b, _ := p.encoder.Marshal(d)
			p.flush(b)
		}
	}
}

func NewPeriodicSnapshotTaker(
	p time.Duration,
	s getter,
	f flusher,
	e marshaler,
) *Periodic {
	return &Periodic{
		period:  p,
		store:   s,
		flush:   f,
		encoder: e,
	}
}
