package vacuum

import "context"

var _ Vacuumer = (*Noop)(nil)

type Noop struct{}

func (n *Noop) Run(context.Context) error { return nil }

func NewNoop() *Noop {
	return &Noop{}
}
