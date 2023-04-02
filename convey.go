package minifac

import "github.com/mazzegi/minifac/grid"

var _ ProducerConsumer = &Conveyor{}

func NewConveyor(name string, dir grid.Direction, capa int) *Conveyor {
	return &Conveyor{
		name:     name,
		dir:      dir,
		capacity: capa,
		buffer:   NewQueue[Resource](),
	}
}

type Conveyor struct {
	name     string
	dir      grid.Direction
	capacity int
	buffer   *Queue[Resource]
}

func (c *Conveyor) Size() grid.Size {
	return grid.S(1, 1)
}

func (c *Conveyor) Dir() grid.Direction {
	return c.dir
}

func (c *Conveyor) ProduceAtPositions(base grid.Position) []grid.Position {
	var pos grid.Position
	switch c.dir {
	case grid.North:
		pos = grid.P(base.X, base.Y-1)
	case grid.East:
		pos = grid.P(base.X+1, base.Y)
	case grid.South:
		pos = grid.P(base.X, base.Y+1)
	default: //West
		pos = grid.P(base.X-1, base.Y)
	}
	return []grid.Position{pos}
}

func (c *Conveyor) ConsumeAtPositions(base grid.Position) []grid.Position {
	//return base.Neighbours()
	return []grid.Position{base}
}

func (c *Conveyor) Tick() {

}

func (c *Conveyor) Name() string {
	return c.name
}

func (c *Conveyor) Consume(res Resource) {
	if !c.CanConsume(res) {
		return
	}
	c.buffer.Enqueue(res)
	Log("%s: consume: %s: buffer=%d/%d", c.name, res, c.buffer.Len(), c.capacity)
}

func (c *Conveyor) CanConsume(Resource) bool {
	return c.buffer.Len() < c.capacity
}

func (c *Conveyor) CanConsumeAny() bool {
	return c.buffer.Len() < c.capacity
}

func (c *Conveyor) Produce() (Resource, bool) {
	res, ok := c.buffer.Dequeue()
	Log("%s: produce: %s: buffer=%d/%d", c.name, res, c.buffer.Len(), c.capacity)
	return res, ok
}

func (c *Conveyor) CanProduce() bool {
	return c.buffer.Len() > 0
}

func (c *Conveyor) Resource() Resource {
	res, ok := c.buffer.Peek()
	if !ok {
		return None
	}
	return res
}
