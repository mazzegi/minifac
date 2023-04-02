package minifac

type Receipt struct {
	input          map[Resource]int
	output         Resource
	productionTime int
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
