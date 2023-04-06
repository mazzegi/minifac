package eeui

import (
	"image"
	"image/color"
	"math"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func NewTextBox(text string, evts *EventHandler) *TextBox {
	b := &TextBox{
		text: text,
	}
	return b
}

type TextBox struct {
	text string
	rect image.Rectangle
}

func (b *TextBox) ChangeText(text string) {
	b.text = text
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
	var twidth fixed.Int26_6
	for _, r := range c.text {
		a, _ := face.GlyphAdvance(r)
		twidth += a
	}
	width2 := int(math.Ceil(float64(twidth) / (64)))
	height2 := int(math.Ceil(float64(theight) / (64)))
	_, _ = width2, height2

	x := c.rect.Min.X + 4
	y := c.rect.Min.Y + 4

	pt := freetype.Pt(x, y)
	_, err := fctx.DrawString(c.text, pt)
	if err != nil {
		panic(err)
	}

}
