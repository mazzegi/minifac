package minifac

import (
	"fmt"
	"sort"

	"github.com/mazzegi/minifac/grid"
	"golang.org/x/exp/maps"
)

var _ ProducerConsumer = &Assembler{}

func NewAssembler(name string, receipt Receipt, inCapa int, outCapa int) *Assembler {
	a := &Assembler{
		name:     name,
		receipt:  receipt,
		inStocks: make(map[Resource]*Stock),
		outStock: NewStock(outCapa),
	}
	for inRes := range receipt.Input {
		a.inStocks[inRes] = NewStock(inCapa)
	}
	return a
}

type Assembler struct {
	name         string
	receipt      Receipt
	inStocks     map[Resource]*Stock
	outStock     *Stock
	lastProdTick int
	currTick     int
	producing    bool
}

func (c *Assembler) Size() grid.Size {
	return grid.S(1, 1)
}

func (c *Assembler) Name() string {
	return c.name
}

func (c *Assembler) Info() []string {
	stockRess := maps.Keys(c.inStocks)
	sort.Slice(stockRess, func(i, j int) bool { return stockRess[i] < stockRess[j] })
	var stock []string
	for _, res := range stockRess {
		stock = append(stock, fmt.Sprintf("Stock: %s: %d", res, c.inStocks[res].TotalAmount()))
	}

	info := []string{
		fmt.Sprintf("Assembler: %s", c.name),
		fmt.Sprintf("Receipt: %s", c.receipt.String()),
	}
	info = append(info, stock...)
	return info
}

func (c *Assembler) ProduceAtPositions(base grid.Position) []grid.Position {
	return base.Neighbours()
}

func (c *Assembler) ConsumeAtPositions(base grid.Position) []grid.Position {
	return []grid.Position{base}
}

func (c *Assembler) Tick() {
	c.currTick++
	if c.producing {
		if c.currTick-c.lastProdTick >= c.receipt.ProductionTime {
			c.outStock.Add(c.receipt.Output, 1)
			c.producing = false
		}
	} else if c.outStock.CanAdd(c.receipt.Output, 1) {
		// see if we can produce something
		for res, cnt := range c.receipt.Input {
			if c.inStocks[res].Amount(res) < cnt {
				return
			}
		}
		// we have enough - take it from stocks
		for res, cnt := range c.receipt.Input {
			c.inStocks[res].Take(res, cnt)
		}
		//start production
		c.producing = true
		c.lastProdTick = c.currTick
	}
}

func (c *Assembler) ConsumeFrom(res Resource, dir grid.Direction) {
	if !c.CanConsumeFrom(res, dir) {
		return
	}
	c.inStocks[res].Add(res, 1)
}

func (c *Assembler) CanConsumeFrom(res Resource, dir grid.Direction) bool {
	_, ok := c.receipt.Input[res]
	if !ok {
		return false
	}
	return c.inStocks[res].CanAdd(res, 1)
}

func (c *Assembler) CanConsumeAny() bool {
	return true
}

func (c *Assembler) Produce() (Resource, bool) {
	if c.outStock.Amount(c.receipt.Output) > 0 {
		//Log("%s: produce: %s", c.name, c.receipt.output)
		c.outStock.Take(c.receipt.Output, 1)
		return c.receipt.Output, true
	}
	return None, false
}

func (c *Assembler) CanProduce() bool {
	return c.outStock.Amount(c.receipt.Output) > 0
}

func (c *Assembler) Resource() Resource {
	return c.receipt.Output
}
