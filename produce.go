package minifac

import "github.com/mazzegi/minifac/grid"

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
	cnt := p.rate.Count(int(p.currTick) - int(p.lastProdTick))
	if cnt > 0 {
		p.stock.Add(p.resource, cnt)
		p.lastProdTick = p.currTick
		Log("%s: tick: %s: stock=%d/%d", p.name, p.resource, p.stock.TotalAmount(), p.stock.capacity)
	}
}

func (p *IncarnationProducer) Size() grid.Size {
	return grid.S(1, 1)
}

func (p *IncarnationProducer) Name() string {
	return p.name
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
		Log("%s: produce: %s: stock=%d/%d", p.name, p.resource, p.stock.TotalAmount(), p.stock.capacity)
		return p.resource, true
	}
	return None, false
}

func (p *IncarnationProducer) Resource() Resource {
	return p.resource
}
