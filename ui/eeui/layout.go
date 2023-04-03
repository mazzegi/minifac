package eeui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func VSplitRect(r image.Rectangle, n int) []image.Rectangle {
	if n <= 0 {
		return []image.Rectangle{}
	}
	hi := r.Dy() / n
	yo := 0
	rs := make([]image.Rectangle, n)
	for i := 0; i < n; i++ {
		if i == n-1 {
			hi = r.Dy() - yo
		}
		//rs[i] = image.Rectangle{r.X, yo, r.W, hi}
		rs[i] = image.Rect(r.Min.X, yo, r.Min.X+r.Dx(), yo+hi)
		yo += hi
	}
	return rs
}

func HSplitRect(r image.Rectangle, n int) []image.Rectangle {
	if n <= 0 {
		return []image.Rectangle{}
	}
	wi := r.Dx() / n
	xo := 0
	rs := make([]image.Rectangle, n)
	for i := 0; i < n; i++ {
		if i == n-1 {
			wi = r.Dx() - xo
		}
		rs[i] = image.Rect(xo, r.Min.Y, xo+wi, r.Min.Y+r.Dy())
		xo += wi
	}
	return rs
}

func NewVBoxLayout(ws ...Widget) *VBoxLayout {
	return &VBoxLayout{
		widgets: ws,
	}
}

type VBoxLayout struct {
	widgets []Widget
	rect    image.Rectangle
}

func (c *VBoxLayout) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
	rs := VSplitRect(c.rect, len(c.widgets))
	for i, w := range c.widgets {
		wr := rs[i]
		w.Resize(&ResizeContext{
			Rect: image.Rect(wr.Min.X, wr.Min.Y, wr.Min.X+wr.Dx(), wr.Min.Y+wr.Dy()),
		})
	}
}

func (c *VBoxLayout) Draw(ctx *DrawContext) {
	screen := ctx.Screen
	rs := VSplitRect(c.rect, len(c.widgets))
	for i, w := range c.widgets {
		wr := rs[i]
		wimg := ebiten.NewImage(wr.Dx(), wr.Dy())
		w.Draw(&DrawContext{
			Screen: wimg,
			Font:   ctx.Font,
		})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, float64(wr.Min.Y))
		screen.DrawImage(wimg, opts)
	}
}

// func (c *VBoxLayout) Resize(x, y, width, height int) {
// 	c.rect = image.Rect(x, y, x+width, y+height)
// 	rs := VSplitRect(c.rect, len(c.widgets))
// 	for i, w := range c.widgets {
// 		wr := rs[i]
// 		w.Resize(wr.Min.X, wr.Min.Y, wr.Dx(), wr.Dy())
// 	}
// }

// func (c *VBoxLayout) Draw(screen *ebiten.Image) {
// 	rs := VSplitRect(c.rect, len(c.widgets))
// 	for i, w := range c.widgets {
// 		wr := rs[i]
// 		wimg := ebiten.NewImage(wr.Dx(), wr.Dy())
// 		w.Draw(wimg)
// 		opts := &ebiten.DrawImageOptions{}
// 		opts.GeoM.Translate(0, float64(wr.Min.Y))
// 		screen.DrawImage(wimg, opts)
// 	}
// }
