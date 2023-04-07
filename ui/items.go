package ui

import (
	"fmt"

	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
)

func CreateObject(ty ImageType, res minifac.Resource) (minifac.Object, error) {
	switch ty {
	case ImageTypeConveyor_east:
		return minifac.NewConveyor("conv", grid.East, 1), nil
	case ImageTypeConveyor_south:
		return minifac.NewConveyor("conv", grid.South, 1), nil
	case ImageTypeConveyor_west:
		return minifac.NewConveyor("conv", grid.West, 1), nil
	case ImageTypeConveyor_north:
		return minifac.NewConveyor("conv", grid.North, 1), nil
	case ImageTypeProducer:
		return minifac.NewIncarnationProducer(string(res)+"_prod", res, minifac.NewRate(1, 2), 2), nil
	case ImageTypeAssembler:
		rec, ok := minifac.ReceiptFor(res)
		if !ok {
			return nil, fmt.Errorf("no receipt for %q", res)
		}
		return minifac.NewAssembler(string(res)+"_ass", rec, 5, 5), nil
	default:
		return nil, fmt.Errorf("invalid item type %q", ty)
	}
}
