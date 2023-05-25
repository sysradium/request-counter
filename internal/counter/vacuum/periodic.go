package vacuum

import (
	"context"
	"time"
)

var _ Vacuumer = (*Periodic)(nil)

type pruner interface {
	Prune()
}

type Periodic struct {
	store  pruner
	period time.Duration
}

func (p *Periodic) Run(ctx context.Context) error {
	t := time.NewTimer(p.period)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			p.store.Prune()
		}
	}
}

func NewPeriodic(p time.Duration, pr pruner) *Periodic {
	return &Periodic{
		period: p,
		store:  pr,
	}
}
