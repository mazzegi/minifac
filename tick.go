package minifac

func NewRate(count, perTicks int) Rate {
	return Rate{Count: count, PerTicks: perTicks}
}

type Rate struct {
	Count    int `json:"count"`
	PerTicks int `json:"per-ticks"`
}

func (r Rate) CalcCount(ticks int) int {
	if ticks%r.PerTicks != 0 {
		return 0
	}
	return r.Count * ticks / r.PerTicks
}
