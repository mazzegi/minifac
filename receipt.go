package minifac

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type Receipt struct {
	Input          map[Resource]int `json:"input"`
	Output         Resource         `json:"output"`
	ProductionTime int              `json:"production-time"`
}

func (r Receipt) String() string {
	ress := maps.Keys(r.Input)
	sort.Slice(ress, func(i, j int) bool { return ress[i] < ress[j] })
	var in []string
	for _, res := range ress {
		in = append(in, fmt.Sprintf("%d %s", r.Input[res], res))
	}
	return fmt.Sprintf("%s -> %s: %d", strings.Join(in, " + "), r.Output, r.ProductionTime)
}

func AllReceipts() []Receipt {
	return []Receipt{
		ReceiptIron(),
		ReceiptSteel(),
	}
}

func ReceiptFor(res Resource) (Receipt, bool) {
	for _, rec := range AllReceipts() {
		if rec.Output == res {
			return rec, true
		}
	}
	return Receipt{}, false
}

func ReceiptIron() Receipt {
	return Receipt{
		Input: map[Resource]int{
			Coal:    2,
			IronOre: 1,
		},
		Output:         Iron,
		ProductionTime: 2,
	}
}

func ReceiptSteel() Receipt {
	return Receipt{
		Input: map[Resource]int{
			Coal: 2,
			Iron: 1,
		},
		Output:         Steel,
		ProductionTime: 3,
	}
}
