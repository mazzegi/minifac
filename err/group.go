package err

import (
	"errors"
)

func NewGroup() *Group {
	return &Group{}
}

type Group struct {
	errs []error
}

func (g *Group) Error() error {
	return errors.Join(g.errs...)
}

func (g *Group) Handle(err error) error {
	if err != nil {
		g.errs = append(g.errs, err)
	}
	return g.Error()
}

func (g *Group) HandleFncs(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			g.errs = append(g.errs, err)
		}
	}
	return g.Error()
}
