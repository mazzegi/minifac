package minifac

func NewStock(capa int) *Stock {
	return &Stock{
		total:     0,
		capacity:  capa,
		resources: make(map[Resource]int),
	}
}

type Stock struct {
	resources map[Resource]int
	total     int
	capacity  int
}

func (s *Stock) Add(res Resource, n int) (added int) {
	add := Min(s.capacity-s.total, n)
	s.total += add
	s.resources[res] += add
	return add
}

func (s *Stock) CanAdd(res Resource, n int) bool {
	return s.total+n <= s.capacity
}

func (s *Stock) Take(res Resource, n int) (taken int) {
	take := Min(Min(s.total, n), s.resources[res])
	s.total -= take
	s.resources[res] -= take
	return take
}

func (s *Stock) Amount(res Resource) int {
	return s.resources[res]
}

func (s *Stock) TotalAmount() int {
	return s.total
}
