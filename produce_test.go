package minifac

import (
	"fmt"
	"testing"
)

func repeat(fn func(), n int) {
	if n <= 0 {
		return
	}
	for i := 0; i < n; i++ {
		fn()
	}
}

func TestProduce(t *testing.T) {
	tests := []struct {
		stock    *Stock
		resource Resource
		rate     Rate
		ticks    int
		expStock int
	}{
		{
			stock:    NewStock(5),
			resource: "Wood",
			rate:     NewRate(2, 1),
			ticks:    2,
			expStock: 4,
		},
		{
			stock:    NewStock(5),
			resource: "Wood",
			rate:     NewRate(2, 1),
			ticks:    3,
			expStock: 5,
		},
		{
			stock:    NewStock(5),
			resource: "Wood",
			rate:     NewRate(2, 1),
			ticks:    10,
			expStock: 5,
		},
		{
			stock:    NewStock(5),
			resource: "Wood",
			rate:     NewRate(1, 3),
			ticks:    10,
			expStock: 3,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_#%02d", i), func(t *testing.T) {
			p := IncarnationProducer{
				resource: test.resource,
				rate:     test.rate,
				stock:    test.stock,
			}
			repeat(p.Tick, test.ticks)
			if test.expStock != p.stock.TotalAmount() {
				t.Fatalf("total: want %d, have %d", test.expStock, p.stock.TotalAmount())
			}
			if test.expStock != p.stock.Amount(test.resource) {
				t.Fatalf("resource: want %d, have %d", test.expStock, p.stock.Amount(test.resource))
			}
		})
	}
}
