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

var (
	ButtonColorNormal = color.RGBA{160, 160, 160, 255}
	ButtonColorHover  = color.RGBA{128, 128, 128, 255}
)

type ButtonState byte

const (
	ButtonStateNormal = iota
	ButtonStateHover
)

func NewButton(text string, evts *EventHandler) *Button {
	b := &Button{
		text:   text,
		margin: 8,
		events: evts,
		state:  ButtonStateNormal,
	}
	evts.OnMouseMove(func(p image.Point) {
		if p.In(b.rect) && b.state != ButtonStateHover {
			b.state = ButtonStateHover
		} else if !p.In(b.rect) && b.state != ButtonStateNormal {
			b.state = ButtonStateNormal
		}
	})
	return b
}

type Button struct {
	text   string
	margin int
	rect   image.Rectangle
	events *EventHandler
	state  ButtonState
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
	w, h := float32(c.rect.Dx()-2*c.margin), float32(c.rect.Dy()-2*c.margin)

	var cr color.Color
	switch c.state {
	case ButtonStateHover:
		cr = ButtonColorHover
	default:
		cr = ButtonColorNormal
	}

	vector.DrawFilledRect(screen, float32(c.margin), float32(c.margin), w, h, cr, true)
	c.drawText(ctx)
}

func (c *Button) drawText(ctx *DrawContext) {
	fontSize := 12
	dpi := 96.0
	fctx := freetype.NewContext()
	fctx.SetDPI(dpi)
	fctx.SetFont(ctx.Font)
	fctx.SetFontSize(float64(fontSize))
	fctx.SetClip(ctx.Screen.Bounds())
	fctx.SetDst(ctx.Screen)
	fctx.SetSrc(image.Black)
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

	x := (ctx.Screen.Bounds().Dx() - int(width2)) / 2
	y := (ctx.Screen.Bounds().Dy() - height2) / 2

	pt := freetype.Pt(x, y)
	_, err := fctx.DrawString(c.text, pt)
	if err != nil {
		panic(err)
	}

}
