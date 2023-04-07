package minifac

import (
	"fmt"

	"github.com/mazzegi/minifac/grid"
)

var _ Consumer = &Finalizer{}

func NewFinalizer(name string, res Resource) *Finalizer {
	return &Finalizer{
		name:     name,
		resource: res,
	}
}

type Finalizer struct {
	name     string
	resource Resource
	total    int
}

func (c *Finalizer) Size() grid.Size {
	return grid.S(1, 1)
}

func (c *Finalizer) Resource() Resource {
	return c.resource
}

func (c *Finalizer) ConsumeAtPositions(base grid.Position) []grid.Position {
	return []grid.Position{base}
}

func (c *Finalizer) Tick() {

}

func (c *Finalizer) Name() string {
	return c.name
}

func (c *Finalizer) Info() []string {
	return []string{
		fmt.Sprintf("Finalizer: %s: %s", c.name, c.resource),
		fmt.Sprintf("Consumed: %d", c.total),
	}
}

func (c *Finalizer) ConsumeFrom(res Resource, dir grid.Direction) {
	if !c.CanConsumeFrom(res, dir) {
		return
	}
	c.total++
}

func (c *Finalizer) CanConsumeFrom(res Resource, dir grid.Direction) bool {
	return res == c.resource
}

func (c *Finalizer) CanConsumeAny() bool {
	return true
}
