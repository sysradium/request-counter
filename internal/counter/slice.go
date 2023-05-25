package counter

import (
	"context"
	"sync"
	"time"

	"github.com/sysradium/request-counter/internal/common"
	"github.com/sysradium/request-counter/internal/counter/vacuum"
)

type SlidingWindowStorage struct {
	m      sync.Mutex
	data   []time.Time
	window time.Duration
	now    func() time.Time
	log    common.Logger
	pruner vacuum.Vacuumer
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *SlidingWindowStorage) Add(t time.Time) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.data = append(s.data, t)

	return nil
}

func (s *SlidingWindowStorage) Len() int {
	s.m.Lock()
	defer s.m.Unlock()

	count := len(s.data)

	cutoffTime := s.now().Add(-s.window)
	for i := 0; i < len(s.data); i++ {
		if s.data[i].After(cutoffTime) {
			break
		}
		count--
	}

	return count
}

func (s *SlidingWindowStorage) Get() []time.Time {
	s.m.Lock()
	defer s.m.Unlock()

	r := make([]time.Time, len(s.data))
	copy(r, s.data)

	return r
}

func (s *SlidingWindowStorage) Start() error {
	if err := s.pruner.Run(s.ctx); err != nil {
		s.log.Printf("unable to start: %v", err)
		return err
	}

	return nil
}

func (s *SlidingWindowStorage) Stop() {
	s.cancel()
}

func (s *SlidingWindowStorage) Prune() {
	cutoffTime := s.now().Add(-s.window)
	s.prune(cutoffTime)
}

func (s *SlidingWindowStorage) prune(cutOffTime time.Time) {
	s.m.Lock()
	defer s.m.Unlock()

	i := 0
	for ; i < len(s.data); i++ {
		if s.data[i].After(cutOffTime) {
			break
		}
	}

	s.log.Printf("prunning %d records ...", i)
	s.data = s.data[i:]
}

func New(
	window time.Duration,
	opts ...Option,
) *SlidingWindowStorage {
	s := &SlidingWindowStorage{
		window: window,
		log:    &common.NullLogger{},
		ctx:    context.Background(),
		now:    time.Now,
		pruner: vacuum.NewNoop(),
	}

	for _, o := range opts {
		o(s)
	}

	s.ctx, s.cancel = context.WithCancel(s.ctx)

	return s
}
