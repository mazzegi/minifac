package ui

import (
	"fmt"

	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
)

type ItemType string

const (
	ItemTypeProducer       ItemType = "producer"
	ItemTypeAssembler      ItemType = "assembler"
	ItemTypeTrash          ItemType = "trash"
	ItemTypeFinalizer      ItemType = "finalizer"
	ItemTypeConveyor_east  ItemType = "conveyor_east"
	ItemTypeConveyor_north ItemType = "conveyor_north"
	ItemTypeConveyor_south ItemType = "conveyor_south"
	ItemTypeConveyor_west  ItemType = "conveyor_west"
)

func (ui *UI) CreateObject(ty ItemType, res minifac.Resource) (minifac.Object, error) {
	switch ty {
	case ItemTypeConveyor_east:
		return minifac.NewConveyor("conv", grid.East, 1), nil
	case ItemTypeConveyor_south:
		return minifac.NewConveyor("conv", grid.South, 1), nil
	case ItemTypeConveyor_west:
		return minifac.NewConveyor("conv", grid.West, 1), nil
	case ItemTypeConveyor_north:
		return minifac.NewConveyor("conv", grid.North, 1), nil
	case ItemTypeProducer:
		conf := ui.config.Producers[res]
		return minifac.NewIncarnationProducer("prod", res, conf.Rate, conf.StockCapacity), nil
	case ItemTypeAssembler:
		conf := ui.config.Assemblers[res]
		rec := minifac.Receipt{
			Input:          conf.Input,
			Output:         res,
			ProductionTime: conf.ProductionTime,
		}
		return minifac.NewAssembler("ass", rec, conf.InputStockCapacity, conf.OutputStockCapacity), nil
	case ItemTypeTrash:
		return minifac.NewTrashbin("trash"), nil
	case ItemTypeFinalizer:
		return minifac.NewFinalizer("finalizer", res), nil
	default:
		return nil, fmt.Errorf("invalid item type %q", ty)
	}
}
