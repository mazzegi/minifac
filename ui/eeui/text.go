package eeui

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
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
	textFunc     func() []string
	rect         image.Rectangle
	lastTextHash string
	img          *ebiten.Image
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
	img := c.createImage(ctx.Font)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(img, opts)
}

func (c *TextBox) textHash() string {
	return strings.Join(c.textFunc(), ":")
}

func (c *TextBox) createImage(bfont *truetype.Font) *ebiten.Image {
	if c.img != nil && c.textHash() == c.lastTextHash {
		return c.img
	}

	img := ebiten.NewImage(c.rect.Dx(), c.rect.Dy())
	vector.DrawFilledRect(img, 0, 0, float32(c.rect.Dx()), float32(c.rect.Dy()), color.Black, true)

	fontSize := 12
	dpi := 96.0
	fctx := freetype.NewContext()
	fctx.SetDPI(dpi)
	fctx.SetFont(bfont)
	fctx.SetFontSize(float64(fontSize))
	fctx.SetClip(img.Bounds())
	fctx.SetDst(img)
	fctx.SetSrc(image.White)
	fctx.SetHinting(font.HintingNone)

	face := truetype.NewFace(bfont, &truetype.Options{
		Size: float64(fontSize),
		DPI:  dpi,
	})
	metrics := face.Metrics()
	_ = metrics
	theight := face.Metrics().Height
	height2 := int(math.Ceil(float64(theight) / (64)))

	x := 4
	y := 24
	for i, text := range c.textFunc() {
		pt := freetype.Pt(x, y+i*height2)
		_, err := fctx.DrawString(text, pt)
		if err != nil {
			panic(err)
		}
	}
	c.lastTextHash = c.textHash()
	c.img = img
	return c.img
}
