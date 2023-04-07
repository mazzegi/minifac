package eeui

import (
	"image"
	"image/color"
	"math"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	ButtonColorNormal = color.RGBA{160, 160, 160, 255}
	//ButtonColorHover  = color.RGBA{128, 128, 128, 255}
)

type ButtonState byte

const (
	ButtonStateNormal = iota
	ButtonStateHover
)

func NewButton(text string, evts *EventHandler) *Button {
	b := &Button{
		text:   text,
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
	rect   image.Rectangle
	events *EventHandler
	state  ButtonState
	img    *ebiten.Image
}

func (b *Button) OnClick(fn func()) {
	b.events.OnMouseLeftClicked(func(p image.Point) {
		if p.In(b.rect) {
			fn()
		}
	})
}

func (b *Button) ChangeText(text string) {
	b.text = text
	b.img = nil
}

func (c *Button) SizeHint() SizeHint {
	return SizeHint{
		MaxHeight: 64,
	}
}

func (c *Button) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
}

func (c *Button) Draw(ctx *DrawContext) {
	screen := ctx.Screen
	if c.img == nil {
		c.createImage(ctx.Font)
	}
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(c.img, opts)

	// var cr color.RGBA
	// switch c.state {
	// case ButtonStateHover:
	// 	cr = ButtonColorNormal
	// default:
	// 	cr = ButtonColorNormal
	// }

	// x, y := float32(c.rect.Min.X), float32(c.rect.Min.Y)
	// w, h := float32(c.rect.Dx()), float32(c.rect.Dy())
	// vector.DrawFilledRect(screen, x, y, w, h, cr, true)
	// c.drawText(ctx)
}

func (c *Button) createImage(bfont *truetype.Font) {
	img := ebiten.NewImage(c.rect.Dx(), c.rect.Dy())
	vector.DrawFilledRect(img, 0, 0, float32(c.rect.Dx()), float32(c.rect.Dy()), ButtonColorNormal, true)

	fontSize := 12
	dpi := 96.0
	fctx := freetype.NewContext()
	fctx.SetDPI(dpi)
	fctx.SetFont(bfont)
	fctx.SetFontSize(float64(fontSize))
	fctx.SetClip(img.Bounds())
	fctx.SetDst(img)
	fctx.SetSrc(image.Black)
	fctx.SetHinting(font.HintingNone)

	face := truetype.NewFace(bfont, &truetype.Options{
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

	x := (c.rect.Dx() - int(width2)) / 2
	y := height2*3/4 + (c.rect.Dy()-height2)/2

	pt := freetype.Pt(x, y)
	_, err := fctx.DrawString(c.text, pt)
	if err != nil {
		panic(err)
	}
	c.img = img
}

// func (c *Button) drawText(ctx *DrawContext) {
// 	fontSize := 12
// 	dpi := 96.0
// 	fctx := freetype.NewContext()
// 	fctx.SetDPI(dpi)
// 	fctx.SetFont(ctx.Font)
// 	fctx.SetFontSize(float64(fontSize))
// 	fctx.SetClip(ctx.Screen.Bounds())
// 	fctx.SetDst(ctx.Screen)
// 	fctx.SetSrc(image.Black)
// 	fctx.SetHinting(font.HintingNone)

// 	face := truetype.NewFace(ctx.Font, &truetype.Options{
// 		Size: float64(fontSize),
// 		DPI:  dpi,
// 	})
// 	metrics := face.Metrics()
// 	_ = metrics
// 	theight := face.Metrics().Height
// 	var twidth fixed.Int26_6
// 	for _, r := range c.text {
// 		a, _ := face.GlyphAdvance(r)
// 		twidth += a
// 	}
// 	width2 := int(math.Ceil(float64(twidth) / (64)))
// 	height2 := int(math.Ceil(float64(theight) / (64)))

// 	x := c.rect.Min.X + (c.rect.Dx()-int(width2))/2
// 	y := c.rect.Min.Y + height2*3/4 + (c.rect.Dy()-height2)/2

// 	pt := freetype.Pt(x, y)
// 	_, err := fctx.DrawString(c.text, pt)
// 	if err != nil {
// 		panic(err)
// 	}

// }
