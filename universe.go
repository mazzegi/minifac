package minifac

import (
	"github.com/mazzegi/minifac/grid"
	"golang.org/x/exp/slices"
)

type Producer interface {
	Produce() (Resource, bool)
	Resource() Resource
	CanProduce() bool
	ProduceAtPositions(base grid.Position) []grid.Position
	Name() string
}

type Consumer interface {
	ConsumeFrom(Resource, grid.Direction)
	CanConsumeFrom(Resource, grid.Direction) bool
	CanConsumeAny() bool
	ConsumeAtPositions(base grid.Position) []grid.Position
	Name() string
}

type ProducerConsumer interface {
	Producer
	Consumer
}

type Object interface {
	Size() grid.Size
	Tick()
	Name() string
	Info() []string
}

func NewUniverse(size grid.Size) *Universe {
	u := &Universe{
		grid: grid.New[Object](size),
	}
	return u
}

type Universe struct {
	grid *grid.Grid[Object]
}

func (u *Universe) Size() grid.Size {
	return u.grid.Size()
}

func (u *Universe) ContainsPosition(p grid.Position) bool {
	return u.grid.ContainsPosition(p)
}

func (u *Universe) AddObject(o Object, at grid.Position) error {
	r := grid.R(at, o.Size())
	return u.grid.Add(o, r)
}

func (u *Universe) DeleteAt(p grid.Position) {
	u.grid.DeleteAt(p)
}

func (u *Universe) AllObjects() []*grid.Object[Object] {
	return u.grid.Objects()
}

func (u *Universe) ObjectAt(p grid.Position) (*grid.Object[Object], bool) {
	o := u.grid.ObjectAt(p)
	if o == nil {
		return nil, false
	}
	return o, true
}

func (u *Universe) Tick() {
	var prods []*grid.Object[Producer]
	var cons []*grid.Object[Consumer]

	findConsumerForPositions := func(fromPos grid.Position, poss []grid.Position, res Resource) (*grid.Object[Consumer], grid.Position, bool) {
		for _, conObj := range cons {
			Log("find-consumer: test %s -> %s", res, conObj.Value.Name())
			for _, pos := range poss {
				fromDir := grid.DirectionFrom(conObj.Position, fromPos)
				if !conObj.Value.CanConsumeFrom(res, fromDir) {
					continue
				}
				if slices.Contains(conObj.Value.ConsumeAtPositions(conObj.Position), pos) {
					return conObj, pos, true
				}
			}
		}
		return nil, grid.Position{}, false
	}

	for _, obj := range u.grid.Objects() {
		obj.Value.Tick()
		if prod, ok := obj.Value.(Producer); ok && prod.CanProduce() {
			prods = append(prods, &grid.Object[Producer]{Value: prod, Rectangle: obj.Rectangle})
		}
		if con, ok := obj.Value.(Consumer); ok && con.CanConsumeAny() {
			cons = append(cons, &grid.Object[Consumer]{Value: con, Rectangle: obj.Rectangle})
		}
	}
	// transport
	for _, prodObj := range prods {
		prodPoss := prodObj.Value.ProduceAtPositions(prodObj.Position)
		res := prodObj.Value.Resource()

		conObj, pos, ok := findConsumerForPositions(prodObj.Position, prodPoss, res)
		if !ok {
			continue
		}
		fromDir := grid.DirectionFrom(conObj.Position, pos)
		prodObj.Value.Produce()
		conObj.Value.ConsumeFrom(res, fromDir)
	}
}
