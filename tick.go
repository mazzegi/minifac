package minifac

func NewRate(count, perTicks int) Rate {
	return Rate{count: count, perTicks: perTicks}
}

type Rate struct {
	count    int
	perTicks int
}

func (r Rate) Count(ticks int) int {
	if ticks%r.perTicks != 0 {
		return 0
	}
	return r.count * ticks / r.perTicks
}
