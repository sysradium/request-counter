package vacuum

import "context"

type Vacuumer interface {
	Run(context.Context) error
}
