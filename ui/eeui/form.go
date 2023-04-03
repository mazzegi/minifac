package eeui

import (
	"image"

	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
)

type SizeHint struct {
	MaxWidth  int
	MaxHeight int
}

type Widget interface {
	Draw(ctx *DrawContext)
	Resize(ctx *ResizeContext)
	SizeHint() SizeHint
}

func NewForm(widget Widget, evts *EventHandler, font *truetype.Font) *Form {
	f := &Form{
		events: evts,
		font:   font,
		widget: widget,
	}
	return f
}

type Form struct {
	rect   image.Rectangle
	events *EventHandler
	font   *truetype.Font
	widget Widget
}

func (f *Form) Resize(x, y, width, height int) {
	f.rect = image.Rect(x, y, x+width, y+height)
	f.widget.Resize(&ResizeContext{Rect: f.rect})
}

func (f *Form) Draw(screen *ebiten.Image) {
	f.widget.Draw(&DrawContext{
		Screen: screen,
		Font:   f.font,
	})
}

type ResizeContext struct {
	Rect image.Rectangle
}

type DrawContext struct {
	Screen *ebiten.Image
	Font   *truetype.Font
}
