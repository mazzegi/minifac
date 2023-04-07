package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
	"github.com/mazzegi/minifac/ui"
)

func main() {

	uni := setupUniverse()
	mfui := ui.New(uni)

	ebiten.SetWindowSize(1024+ui.MenuWidth, 1024)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(800, 640, -1, -1)
	ebiten.SetWindowTitle("MiniFAC")
	if err := ebiten.RunGame(mfui); err != nil {
		log.Fatal(err)
	}
}

func setupUniverse() *minifac.Universe {
	size := grid.S(16, 16)

	u := minifac.NewUniverse(size)
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
	u.AddObject(minifac.NewConveyor("conv_steel_3", grid.East, 1), grid.P(10, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_4", grid.East, 1), grid.P(11, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_5", grid.East, 1), grid.P(12, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_6", grid.East, 1), grid.P(13, 3))
	u.AddObject(minifac.NewConveyor("conv_steel_7", grid.East, 1), grid.P(14, 3))

	u.AddObject(minifac.NewTrashbin("trashbin_1"), grid.P(15, 3))

	return u
}

func setupUniverse2() *minifac.Universe {
	size := grid.S(16, 16)

	u := minifac.NewUniverse(size)
	u.AddObject(minifac.NewIncarnationProducer("prod_iron_ore", minifac.IronOre, minifac.NewRate(1, 2), 2), grid.P(1, 1))
	u.AddObject(minifac.NewConveyor("conv_ironore_1", grid.East, 1), grid.P(2, 1))

	u.AddObject(minifac.NewAssembler("ass_iron", minifac.ReceiptIron(), 5, 5), grid.P(3, 1))

	u.AddObject(minifac.NewIncarnationProducer("prod_coal", minifac.Coal, minifac.NewRate(1, 2), 2), grid.P(3, 3))
	u.AddObject(minifac.NewConveyor("conv_coal_1", grid.North, 1), grid.P(3, 2))

	return u
}
