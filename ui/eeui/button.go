package eeui

import (
	"image"
	"image/color"

	"github.com/goki/freetype"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

func NewButton(text string, evts *EventHandler) *Button {
	return &Button{
		text:   text,
		margin: 8,
		events: evts,
	}
}

type Button struct {
	text   string
	margin int
	rect   image.Rectangle
	events *EventHandler
}

func (b *Button) OnClick(fn func()) {
	b.events.OnMouseLeftClicked(func(p image.Point) {
		if p.In(b.rect) {
			fn()
		}
	})
}

func (c *Button) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
}

func (c *Button) Draw(ctx *DrawContext) {
	screen := ctx.Screen
	//x0, y0 := c.rect.Min.X+c.margin, c.rect.Min.Y+c.margin
	w, h := float32(c.rect.Dx()-2*c.margin), float32(c.rect.Dy()-2*c.margin)
	vector.DrawFilledRect(screen, float32(c.margin), float32(c.margin), w, h, color.RGBA{128, 128, 128, 255}, true)
	//vector.DrawFilledRect(screen, float32(x0), float32(y0), float32(c.rect.Dx()-2*c.margin), float32(c.rect.Dy()-2*c.margin), color.RGBA{128, 128, 128, 255}, true)
	c.drawText(ctx)
}

func (c *Button) drawText(ctx *DrawContext) {
	fctx := freetype.NewContext()
	fctx.SetDPI(120)
	fctx.SetFont(ctx.Font)
	fctx.SetFontSize(12)
	fctx.SetClip(ctx.Screen.Bounds())
	fctx.SetDst(ctx.Screen)
	fctx.SetSrc(image.Black)
	fctx.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 30)
	_, err := fctx.DrawString(c.text, pt)
	if err != nil {
		panic(err)
	}
}
