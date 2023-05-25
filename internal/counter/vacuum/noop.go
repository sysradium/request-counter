package vacuum

import "context"

var _ Vacuumer = (*Noop)(nil)

// Noop is a mock vacuumer, used if vacuum must be disabled, or for tests
type Noop struct{}

func (n *Noop) Run(context.Context) error { return nil }

func NewNoop() *Noop {
	return &Noop{}
}
