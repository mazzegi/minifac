package ui

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mazzegi/minifac"
)

func NewUI(uni *minifac.Universe) *UI {
	ui := &UI{
		eventHandler: NewHandler(),
		universe:     uni,
		imageHandler: NewImageHandler(uni),
		ticker:       time.NewTicker(100 * time.Millisecond),
	}
	ui.eventHandler.OnMouseLeftClicked(func(p Pos) {
		x, y := p.X/int(ui.scaleX), p.Y/int(ui.scaleY)
		//ui.universe.OnLeftClick(Position{x, y})
		_, _ = x, y
	})

	return ui
}

type UI struct {
	dx, dy         int
	scaleX, scaleY float64
	eventHandler   *Handler
	universe       *minifac.Universe
	imageHandler   *ImageHandler
	ticker         *time.Ticker
}

func (ui *UI) Update() error {
	ui.eventHandler.Update()
	select {
	case <-ui.ticker.C:
		ui.universe.Tick()
	default:
	}
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.2f", ebiten.ActualTPS()))
	pimgs := ui.imageHandler.Images()
	for _, pimg := range pimgs {
		bs := pimg.Image.Bounds()
		scaleX, scaleY := ui.scaleX/float64(bs.Dx()), ui.scaleY/float64(bs.Dy())

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(scaleX, scaleY)
		opts.GeoM.Translate(ui.scaleX*float64(pimg.Position.X), ui.scaleY*float64(pimg.Position.Y))
		screen.DrawImage(pimg.Image, opts)
	}
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ui.dx = outsideWidth
	ui.dy = outsideHeight
	sz := ui.universe.Size()
	ui.scaleX, ui.scaleY = float64(ui.dx)/float64(sz.DX), float64(ui.dy)/float64(sz.DY)
	return outsideWidth, outsideHeight
}
