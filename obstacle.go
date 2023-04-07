package minifac

import (
	"fmt"

	"github.com/mazzegi/minifac/grid"
)

type ObstacleType string

const (
	ObstacleWall ObstacleType = "wall"
)

var _ Object = &Obstacle{}

func NewObstacle(name string, typ ObstacleType) *Obstacle {
	return &Obstacle{
		name: name,
		typ:  typ,
	}
}

type Obstacle struct {
	name string
	typ  ObstacleType
}

func (c *Obstacle) Type() ObstacleType {
	return c.typ
}

func (c *Obstacle) Size() grid.Size {
	return grid.S(1, 1)
}

func (c *Obstacle) Tick() {

}

func (c *Obstacle) Name() string {
	return c.name
}

func (c *Obstacle) Info() []string {
	return []string{
		fmt.Sprintf("Obstacle: %s", c.name),
	}
}
