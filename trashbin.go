package minifac

import "github.com/mazzegi/minifac/grid"

var _ Consumer = &Trashbin{}

func NewTrashbin(name string) *Trashbin {
	return &Trashbin{
		name: name,
	}
}

type Trashbin struct {
	name string
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

func (c *Trashbin) Consume(res Resource) {
	Log("%s: consume: %s", c.name, res)
}

func (c *Trashbin) CanConsume(Resource) bool {
	return true
}

func (c *Trashbin) CanConsumeAny() bool {
	return true
}
