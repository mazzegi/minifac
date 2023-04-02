package main

import (
	"time"

	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
)

func main() {
	u := minifac.NewUniverse()
	u.AddObject(minifac.NewIncarnationProducer("prod_iron", minifac.Iron, minifac.NewRate(1, 2), 2), grid.P(1, 1))
	u.AddObject(minifac.NewConveyor("conv_iron_1", grid.East, 1), grid.P(2, 1))
	u.AddObject(minifac.NewConveyor("conv_iron_2", grid.East, 1), grid.P(3, 1))
	u.AddObject(minifac.NewConveyor("conv_iron_3", grid.East, 1), grid.P(4, 1))
	u.AddObject(minifac.NewConveyor("conv_iron_4", grid.East, 1), grid.P(5, 1))
	u.AddObject(minifac.NewConveyor("conv_iron_5", grid.East, 1), grid.P(6, 1))
	u.AddObject(minifac.NewConveyor("v_conv_iron_1", grid.South, 1), grid.P(7, 1))
	u.AddObject(minifac.NewConveyor("v_conv_iron_2", grid.South, 1), grid.P(7, 2))

	u.AddObject(minifac.NewIncarnationProducer("prod_coal", minifac.Coal, minifac.NewRate(1, 2), 2), grid.P(1, 5))
	u.AddObject(minifac.NewConveyor("conv_coal_1", grid.East, 1), grid.P(2, 5))
	u.AddObject(minifac.NewConveyor("conv_coal_2", grid.East, 1), grid.P(3, 5))
	u.AddObject(minifac.NewConveyor("conv_coal_3", grid.East, 1), grid.P(4, 5))
	u.AddObject(minifac.NewConveyor("conv_coal_4", grid.East, 1), grid.P(5, 5))
	u.AddObject(minifac.NewConveyor("conv_coal_5", grid.East, 1), grid.P(6, 5))
	u.AddObject(minifac.NewConveyor("v_conv_coal_1", grid.North, 1), grid.P(7, 5))
	u.AddObject(minifac.NewConveyor("v_conv_coal_2", grid.North, 1), grid.P(7, 4))

	u.AddObject(minifac.NewAssembler("ass_steel", minifac.ReceiptSteel(), 5, 5), grid.P(7, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_1", grid.East, 1), grid.P(8, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_2", grid.East, 1), grid.P(9, 3))

	u.AddObject(minifac.NewTrashbin("trashbin_1"), grid.P(10, 3))

	ticks := 200
	tickSleep := 500 * time.Millisecond
	for i := 0; i < ticks; i++ {
		minifac.Log("*** tick %02d ***", i+1)
		u.Tick()
		<-time.After(tickSleep)
	}
}
