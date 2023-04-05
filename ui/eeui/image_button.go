package eeui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewImageButton(img *ebiten.Image, dx, dy int, evts *EventHandler) *ImageButton {
	//scale image
	scaledImg := ebiten.NewImage(dx, dy)
	opts := &ebiten.DrawImageOptions{}
	scaleX := float64(dx) / float64(img.Bounds().Dx())
	scaleY := float64(dy) / float64(img.Bounds().Dy())
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	opts.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(img, opts)

	b := &ImageButton{
		events: evts,
		img:    scaledImg,
		dx:     dx,
		dy:     dy,
	}
	evts.OnMouseMove(func(p image.Point) {

	})
	return b
}

type ImageButton struct {
	rect   image.Rectangle
	events *EventHandler
	img    *ebiten.Image
	dx, dy int
}

func (b *ImageButton) OnClick(fn func()) {
	b.events.OnMouseLeftClicked(func(p image.Point) {
		if p.In(b.rect) {
			fn()
		}
	})
}

func (c *ImageButton) SizeHint() SizeHint {
	return SizeHint{
		MaxHeight: c.dy,
	}
}

func (c *ImageButton) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
}

func (c *ImageButton) Draw(ctx *DrawContext) {
	screen := ctx.Screen
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(c.img, opts)
}
