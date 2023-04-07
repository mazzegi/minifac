package minifac

import (
	"fmt"

	"github.com/mazzegi/minifac/grid"
)

var _ Consumer = &Trashbin{}

func NewTrashbin(name string) *Trashbin {
	return &Trashbin{
		name: name,
	}
}

type Trashbin struct {
	name  string
	total int
}

func (c *Trashbin) Size() grid.Size {
	return grid.S(1, 1)
}

func (c *Trashbin) ConsumeAtPositions(base grid.Position) []grid.Position {
	return []grid.Position{base}
}

func (c *Trashbin) Tick() {

}

func (c *Trashbin) Name() string {
	return c.name
}

func (c *Trashbin) Info() []string {
	return []string{
		fmt.Sprintf("Trashbin: %s", c.name),
		fmt.Sprintf("Consumed: %d", c.total),
	}
}

func (c *Trashbin) Consume(res Resource) {
	c.total++
	//Log("%s: consume: %s", c.name, res)
}

func (c *Trashbin) CanConsume(Resource) bool {
	return true
}

func (c *Trashbin) CanConsumeAny() bool {
	return true
}
