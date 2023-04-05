package eeui

import (
	"image"
)

func VSplitRectBySizeHints(r image.Rectangle, sty BoxLayoutStyles, sizeHints []SizeHint) []image.Rectangle {
	n := len(sizeHints)
	if n <= 0 {
		return []image.Rectangle{}
	}
	//apply padding
	r = image.Rect(
		r.Min.X+sty.Padding,
		r.Min.Y+sty.Padding,
		r.Max.X-sty.Padding,
		r.Max.Y-sty.Padding,
	)

	hi := (r.Dy() - (n-1)*sty.Gap) / n
	yo := r.Min.Y
	rs := make([]image.Rectangle, n)
	for i := 0; i < n; i++ {
		if i == n-1 {
			hi = r.Dy() - yo
		}
		sh := sizeHints[i]
		if sh.MaxHeight > 0 && hi > sh.MaxHeight {
			hi = sh.MaxHeight
		}

		rs[i] = image.Rect(r.Min.X, yo, r.Min.X+r.Dx(), yo+hi)
		yo += hi + sty.Gap
	}
	return rs
}

func HSplitRectBySizeHints(r image.Rectangle, sty BoxLayoutStyles, sizeHints []SizeHint) []image.Rectangle {
	n := len(sizeHints)
	if n <= 0 {
		return []image.Rectangle{}
	}
	//apply padding
	r = image.Rect(
		r.Min.X+sty.Padding,
		r.Min.Y+sty.Padding,
		r.Max.X-sty.Padding,
		r.Max.Y-sty.Padding,
	)

	wi := (r.Dx() - (n-1)*sty.Gap) / n
	xo := r.Min.X
	rs := make([]image.Rectangle, n)
	for i := 0; i < n; i++ {
		if i == n-1 {
			wi = r.Max.X - xo
		}
		sh := sizeHints[i]
		if sh.MaxWidth > 0 && wi > sh.MaxWidth {
			wi = sh.MaxWidth
		}
		rs[i] = image.Rect(xo, r.Min.Y, xo+wi, r.Min.Y+r.Dy())
		xo += wi + sty.Gap
	}
	return rs
}

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

type BoxLayoutStyles struct {
	Padding int
	Gap     int
}

func NewVBoxLayout(styles BoxLayoutStyles, ws ...Widget) *VBoxLayout {
	return &VBoxLayout{
		widgets: ws,
		styles:  styles,
	}
}

type VBoxLayout struct {
	widgets []Widget
	styles  BoxLayoutStyles
	rect    image.Rectangle
}

func (c *VBoxLayout) SizeHint() SizeHint {
	return SizeHint{}
}

func (c *VBoxLayout) sizeHints() []SizeHint {
	shs := make([]SizeHint, len(c.widgets))
	for i, w := range c.widgets {
		shs[i] = w.SizeHint()
	}
	return shs
}

func (c *VBoxLayout) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
	rs := VSplitRectBySizeHints(c.rect, c.styles, c.sizeHints())
	for i, w := range c.widgets {
		wr := rs[i]
		w.Resize(&ResizeContext{
			Rect: wr,
		})
	}
}

func (c *VBoxLayout) Draw(ctx *DrawContext) {
	for _, w := range c.widgets {
		w.Draw(ctx)
	}
}

// HBox
func NewHBoxLayout(styles BoxLayoutStyles, ws ...Widget) *HBoxLayout {
	return &HBoxLayout{
		widgets: ws,
		styles:  styles,
	}
}

type HBoxLayout struct {
	widgets []Widget
	styles  BoxLayoutStyles
	rect    image.Rectangle
}

func (c *HBoxLayout) SizeHint() SizeHint {
	return SizeHint{}
}

func (c *HBoxLayout) sizeHints() []SizeHint {
	shs := make([]SizeHint, len(c.widgets))
	for i, w := range c.widgets {
		shs[i] = w.SizeHint()
	}
	return shs
}

func (c *HBoxLayout) Resize(ctx *ResizeContext) {
	c.rect = ctx.Rect
	rs := HSplitRectBySizeHints(c.rect, c.styles, c.sizeHints())
	for i, w := range c.widgets {
		wr := rs[i]
		w.Resize(&ResizeContext{
			Rect: wr,
		})
	}
}

func (c *HBoxLayout) Draw(ctx *DrawContext) {
	for _, w := range c.widgets {
		w.Draw(ctx)
	}
}
