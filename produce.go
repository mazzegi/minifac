package minifac

import (
	"fmt"

	"github.com/mazzegi/minifac/grid"
)

var _ Producer = &IncarnationProducer{}

func NewIncarnationProducer(name string, res Resource, rate Rate, stockCapa int) *IncarnationProducer {
	return &IncarnationProducer{
		name:     name,
		resource: res,
		rate:     rate,
		stock:    NewStock(stockCapa),
	}
}

type IncarnationProducer struct {
	name         string
	resource     Resource
	rate         Rate
	stock        *Stock
	lastProdTick int
	currTick     int
}

func (p *IncarnationProducer) Tick() {
	p.currTick++
	cnt := p.rate.CalcCount(int(p.currTick) - int(p.lastProdTick))
	if cnt > 0 {
		p.stock.Add(p.resource, cnt)
		p.lastProdTick = p.currTick
	}
}

func (p *IncarnationProducer) Size() grid.Size {
	return grid.S(1, 1)
}

func (p *IncarnationProducer) Name() string {
	return p.name
}

func (p *IncarnationProducer) Info() []string {
	return []string{
		fmt.Sprintf("Incarnation Producer: %s", p.name),
		fmt.Sprintf("Resource: %s", p.resource),
		fmt.Sprintf("Rate    : %d/%d", p.rate.Count, p.rate.PerTicks),
		fmt.Sprintf("Stock   : %d", p.stock.TotalAmount()),
	}
}

func (p *IncarnationProducer) CanProduce() bool {
	return p.stock.Amount(p.resource) > 0
}

func (p *IncarnationProducer) ProduceAtPositions(base grid.Position) []grid.Position {
	return base.Neighbours()
}

func (p *IncarnationProducer) Produce() (Resource, bool) {
	if p.stock.Amount(p.resource) > 0 {
		p.stock.Take(p.resource, 1)
		return p.resource, true
	}
	return NoResource, false
}

func (p *IncarnationProducer) Resource() Resource {
	return p.resource
}
