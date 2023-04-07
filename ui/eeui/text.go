package eeui

import (
	"image"
	"image/color"
	"math"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

func NewTextBox(evts *EventHandler) *TextBox {
	b := &TextBox{
		textFunc: func() []string { return []string{} },
	}
	return b
}

type TextBox struct {
	textFunc func() []string
	rect     image.Rectangle
}

func (b *TextBox) ChangeTextFunc(fn func() []string) {
	b.textFunc = fn
}

func (c *TextBox) SizeHint() SizeHint {
	return SizeHint{}
}

func (c *TextBox) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
}

func (c *TextBox) Draw(ctx *DrawContext) {
	screen := ctx.Screen
	x, y := float32(c.rect.Min.X), float32(c.rect.Min.Y)
	w, h := float32(c.rect.Dx()), float32(c.rect.Dy())
	vector.DrawFilledRect(screen, x, y, w, h, color.Black, true)
	c.drawText(ctx)
}

func (c *TextBox) drawText(ctx *DrawContext) {
	fontSize := 12
	dpi := 96.0
	fctx := freetype.NewContext()
	fctx.SetDPI(dpi)
	fctx.SetFont(ctx.Font)
	fctx.SetFontSize(float64(fontSize))
	fctx.SetClip(ctx.Screen.Bounds())
	fctx.SetDst(ctx.Screen)
	fctx.SetSrc(image.White)
	fctx.SetHinting(font.HintingNone)

	face := truetype.NewFace(ctx.Font, &truetype.Options{
		Size: float64(fontSize),
		DPI:  dpi,
	})
	metrics := face.Metrics()
	_ = metrics
	theight := face.Metrics().Height
	height2 := int(math.Ceil(float64(theight) / (64)))

	x := c.rect.Min.X + 4
	y := c.rect.Min.Y + 4

	for i, text := range c.textFunc() {
		pt := freetype.Pt(x, y+i*height2)
		_, err := fctx.DrawString(text, pt)
		if err != nil {
			panic(err)
		}
	}
}
