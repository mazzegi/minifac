package ui

func NewForm(evts EventHandler) *Form {
	f := &Form{
		events: evts,
	}
	return f
}

type Form struct {
	events EventHandler
}
