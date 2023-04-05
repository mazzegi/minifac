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
	Consume(Resource)
	CanConsume(Resource) bool
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

	findConsumerForPositions := func(poss []grid.Position, res Resource) (*grid.Object[Consumer], bool) {
		for _, conObj := range cons {
			if !conObj.Value.CanConsume(res) {
				continue
			}
			for _, pos := range poss {
				if slices.Contains(conObj.Value.ConsumeAtPositions(conObj.Position), pos) {
					return conObj, true
				}
			}
		}
		return nil, false
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

		name := prodObj.Value.Name()
		Log("looking for consumers: %s", name)
		conObj, ok := findConsumerForPositions(prodPoss, res)
		if !ok {
			continue
		}
		prodObj.Value.Produce()
		conObj.Value.Consume(res)
	}
}
