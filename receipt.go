package minifac

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type Receipt struct {
	input          map[Resource]int
	output         Resource
	productionTime int
}

func (r Receipt) String() string {
	ress := maps.Keys(r.input)
	sort.Slice(ress, func(i, j int) bool { return ress[i] < ress[j] })
	var in []string
	for _, res := range ress {
		in = append(in, fmt.Sprintf("%d %s", r.input[res], res))
	}
	return fmt.Sprintf("%s -> %s: %d", strings.Join(in, " + "), r.output, r.productionTime)
}

func ReceiptSteel() Receipt {
	return Receipt{
		input: map[Resource]int{
			Coal: 2,
			Iron: 1,
		},
		output:         Steel,
		productionTime: 3,
	}
}
