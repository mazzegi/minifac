package grid

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"
)

type Direction byte

const (
	North Direction = iota
	East
	South
	West
)

func P(x, y int) Position {
	return Position{x, y}
}

func S(dx, dy int) Size {
	return Size{dx, dy}
}

func R(p Position, s Size) Rectangle {
	return Rectangle{Position: p, Size: s}
}

type Position struct {
	X, Y int
}

func (p Position) Less(q Position) bool {
	switch {
	case p.Y < q.Y:
		return true
	case p.Y > q.Y:
		return false
	default: // p.Y == q.Y
		return p.X < q.X
	}
}

func (p Position) Neighbours() []Position {
	return []Position{
		P(p.X-1, p.Y),
		P(p.X+1, p.Y),
		P(p.X, p.Y-1),
		P(p.X, p.Y+1),
	}
}

type Size struct {
	DX, DY int
}

type Rectangle struct {
	Position
	Size
}

func (p Position) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("%d,%d+%dx%d", r.X, r.Y, r.DX, r.DY)
}

func (r Rectangle) Positions() []Position {
	var poss []Position
	for x := r.X; x < r.X+r.DX; x++ {
		for y := r.Y; y < r.Y+r.DY; y++ {
			poss = append(poss, P(x, y))
		}
	}
	return poss
}

type Object[T any] struct {
	Value T
	Rectangle
}

func New[T any](size Size) *Grid[T] {
	return &Grid[T]{
		size:    size,
		objects: map[Position]*Object[T]{},
	}
}

type Grid[T any] struct {
	size    Size
	objects map[Position]*Object[T]
}

func (g *Grid[T]) Size() Size {
	return g.size
}

func (g *Grid[T]) ContainsPosition(p Position) bool {
	return p.X >= 0 && p.X < g.size.DX &&
		p.Y >= 0 && p.Y < g.size.DY
}

func (g *Grid[T]) CanAddRectangle(r Rectangle) bool {
	for _, p := range r.Positions() {
		if _, occ := g.objects[p]; occ {
			return false
		}
	}
	return true
}

func (g *Grid[T]) Add(t T, r Rectangle) error {
	if !g.CanAddRectangle(r) {
		return fmt.Errorf("rectangle %s is already occupied", r)
	}
	o := &Object[T]{Value: t, Rectangle: r}
	for _, p := range r.Positions() {
		g.objects[p] = o
	}
	return nil
}

func (g *Grid[T]) DeleteAt(p Position) {
	delete(g.objects, p)
}

func (g *Grid[T]) ObjectAt(p Position) *Object[T] {
	return g.objects[p]
}

func (g *Grid[T]) Objects() []*Object[T] {
	poss := maps.Keys(g.objects)
	sort.Slice(poss, func(i, j int) bool {
		return poss[i].Less(poss[j])
	})
	objs := make([]*Object[T], len(poss))
	for i, pos := range poss {
		objs[i] = g.objects[pos]
	}
	return objs
}
