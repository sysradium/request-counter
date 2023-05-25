package counter

import (
	"context"
	"time"

	"github.com/sysradium/request-counter/internal/common"
	"github.com/sysradium/request-counter/internal/counter/vacuum"
)

type Option func(*SlidingWindowStorage)

func WithLogger(l common.Logger) Option {
	return func(s *SlidingWindowStorage) {
		s.log = l
	}
}

func WithContext(ctx context.Context) Option {
	return func(s *SlidingWindowStorage) {
		s.ctx = ctx
	}
}

func WithClock(clock func() time.Time) Option {
	return func(s *SlidingWindowStorage) {
		s.now = clock
	}
}

func WithVacuumer(v vacuum.Vacuumer) Option {
	return func(s *SlidingWindowStorage) {
		s.pruner = v
	}
}

func WithPeriodicVacuum(p time.Duration) Option {
	return func(s *SlidingWindowStorage) {
		s.pruner = vacuum.NewPeriodic(p, s)
	}
}

func WithData(d []time.Time) Option {
	return func(s *SlidingWindowStorage) {
		s.data = d
	}
}
