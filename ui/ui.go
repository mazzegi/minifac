package ui

import (
	"fmt"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
	"github.com/mazzegi/minifac/ui/eeui"
)

const MenuWidth = 360

func New(uni *minifac.Universe) *UI {
	evts := eeui.NewHandler()
	ui := &UI{
		eventHandler: evts,
		universe:     uni,
		imageHandler: NewImageHandler(uni),
		ticker:       time.NewTicker(500 * time.Millisecond),
		running:      false,
	}
	ui.ticker.Stop()

	infoBox := eeui.NewTextBox(evts)
	selectItem := func(ty ImageType, res minifac.Resource) {
		ui.selectedItem = ty
		ui.selectedResource = res
		infoBox.ChangeTextFunc(func() []string {
			return []string{
				"Selected:",
				fmt.Sprintf("Item    : %s", ty),
				fmt.Sprintf("Resource: %s", res),
			}
		})
	}

	baseTickerTime := 500 * time.Millisecond
	tickerTime := baseTickerTime
	resetTicker := func(d time.Duration) {
		tickerTime = d
		infoBox.ChangeTextFunc(func() []string {
			return []string{
				fmt.Sprintf("Ticker: %s", d),
			}
		})
		if ui.running {
			ui.ticker.Reset(d)
		}
	}

	startBtn := eeui.NewButton("start", evts)
	startBtn.OnClick(func() {
		switch {
		case !ui.running:
			ui.running = true
			resetTicker(tickerTime)
			startBtn.ChangeText("stop")
		default: //running
			ui.running = false
			ui.ticker.Stop()
			startBtn.ChangeText("start")
		}
	})
	startLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		startBtn,
	)

	btnBase := eeui.NewButton("Base Speed", evts)
	btnBase.OnClick(func() {
		resetTicker(baseTickerTime)
	})
	btnFaster := eeui.NewButton("Faster", evts)
	btnFaster.OnClick(func() {
		tickerTime -= 50 * time.Millisecond
		if tickerTime < 50*time.Millisecond {
			tickerTime = 50 * time.Millisecond
		}
		resetTicker(tickerTime)
	})
	tickerLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		btnBase, btnFaster,
	)

	btnConvEast := eeui.NewImageButton(mustLoadImage(ImageTypeConveyor_east), 48, 48, evts)
	btnConvEast.OnClick(func() {
		selectItem(ImageTypeConveyor_east, minifac.None)
	})
	btnConvSouth := eeui.NewImageButton(mustLoadImage(ImageTypeConveyor_south), 48, 48, evts)
	btnConvSouth.OnClick(func() {
		selectItem(ImageTypeConveyor_south, minifac.None)
	})
	btnConvWest := eeui.NewImageButton(mustLoadImage(ImageTypeConveyor_west), 48, 48, evts)
	btnConvWest.OnClick(func() {
		selectItem(ImageTypeConveyor_west, minifac.None)
	})
	btnConvNorth := eeui.NewImageButton(mustLoadImage(ImageTypeConveyor_north), 48, 48, evts)
	btnConvNorth.OnClick(func() {
		selectItem(ImageTypeConveyor_north, minifac.None)
	})
	convLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		btnConvEast,
		btnConvSouth,
		btnConvWest,
		btnConvNorth,
	)

	//Production
	var prodBtns []eeui.Widget
	for _, bres := range minifac.BaseResources() {
		bres := bres
		btn := eeui.NewImageButton(ui.imageHandler.createThumbnailOverlay(ImageTypeProducer, resourceImageType(bres)), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ImageTypeProducer, bres)
		})
		prodBtns = append(prodBtns, btn)
	}
	prodLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		prodBtns...,
	)

	//Assembly
	var assBtns []eeui.Widget
	for _, rec := range minifac.AllReceipts() {
		rec := rec
		btn := eeui.NewImageButton(ui.imageHandler.createThumbnailOverlay(ImageTypeAssembler, resourceImageType(rec.Output)), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ImageTypeAssembler, rec.Output)
		})
		assBtns = append(assBtns, btn)
	}
	assLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		assBtns...,
	)

	//Misc
	var miscBtns []eeui.Widget
	{
		btn := eeui.NewImageButton(ui.imageHandler.images[ImageTypeTrash], 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ImageTypeTrash, minifac.None)
		})
		miscBtns = append(miscBtns, btn)
	}
	miscLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		miscBtns...,
	)

	layout := eeui.NewVBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     24,
		},
		startLayout,
		tickerLayout,
		convLayout,
		prodLayout,
		assLayout,
		miscLayout,
		infoBox,
	)

	font := mustLoadFont("fonts/inter/Inter-Medium.ttf")
	menu := eeui.NewForm(layout, evts, font)
	ui.menu = menu

	ui.eventHandler.OnMouseRightClicked(func(p image.Point) {
		x, y := p.X/int(ui.scaleX), p.Y/int(ui.scaleY)
		ui.universe.DeleteAt(grid.P(x, y))
	})
	ui.eventHandler.OnMouseLeftClicked(func(p image.Point) {
		x, y := p.X/int(ui.scaleX), p.Y/int(ui.scaleY)
		pos := grid.P(x, y)
		if !ui.universe.ContainsPosition(pos) {
			return
		}
		exobj, ok := ui.universe.ObjectAt(pos)
		if !ok {
			// add new object
			obj, err := CreateObject(ui.selectedItem, ui.selectedResource)
			if err != nil {
				minifac.Log("ERROR: create-object: %v", err)
				return
			}
			ui.universe.AddObject(obj, pos)
		} else {
			infoBox.ChangeTextFunc(exobj.Value.Info)
		}
	})

	return ui
}

type UI struct {
	dx, dy           int
	scaleX, scaleY   float64
	eventHandler     *eeui.EventHandler
	universe         *minifac.Universe
	imageHandler     *ImageHandler
	ticker           *time.Ticker
	running          bool
	menu             *eeui.Form
	selectedItem     ImageType
	selectedResource minifac.Resource
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

func (ui *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ui.menu.Resize(outsideWidth-MenuWidth, 0, MenuWidth, outsideHeight)

	ui.dx = outsideWidth - MenuWidth
	ui.dy = outsideHeight
	sz := ui.universe.Size()
	ui.scaleX, ui.scaleY = float64(ui.dx)/float64(sz.DX), float64(ui.dy)/float64(sz.DY)
	return outsideWidth, outsideHeight
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
	ui.menu.Draw(screen)
}
