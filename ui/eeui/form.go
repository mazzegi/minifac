package eeui

import "github.com/hajimehoshi/ebiten/v2"

type Layout interface{}

type Widget interface{}

func NewForm(layout Layout, evts EventHandler) *Form {
	f := &Form{
		events: evts,
		layout: layout,
	}
	return f
}

type Form struct {
	dx, dy int
	events EventHandler
	layout Layout
}

func (f *Form) Resize(width, height int) {
	f.dx, f.dy = width, height
}

func (f *Form) Draw(screen *ebiten.Image) {

}
