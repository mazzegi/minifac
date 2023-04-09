package ui

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/assets"
	"github.com/mazzegi/minifac/grid"
	"github.com/mazzegi/minifac/ui/eeui"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const MenuWidth = 360

func New(uni *minifac.Universe, assets *assets.Assets, config *minifac.Config) *UI {
	evts := eeui.NewHandler()
	ui := &UI{
		eventHandler: evts,
		universe:     uni,
		assets:       assets,
		config:       config,
		imageHandler: NewImageHandler(uni, assets),
		ticker:       time.NewTicker(500 * time.Millisecond),
		running:      false,
	}
	ui.ticker.Stop()

	infoBox := eeui.NewTextBox(evts)
	selectItem := func(ty ItemType, res minifac.Resource) {
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

	var convBtns []eeui.Widget
	for _, dir := range []grid.Direction{grid.East, grid.South, grid.West, grid.North} {
		dir := dir
		btn := eeui.NewImageButton(ui.assets.Conveyor(dir), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(conveyorType(dir), minifac.NoResource)
		})
		convBtns = append(convBtns, btn)
	}
	convLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		convBtns...,
	)

	//Production
	var prodBtns []eeui.Widget
	for _, bres := range ui.config.BaseResources {
		bres := bres
		btn := eeui.NewImageButton(ui.imageHandler.createThumbnailOverlay(ItemTypeProducer, bres), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ItemTypeProducer, bres)
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
	assRess := maps.Keys(ui.config.Assemblers)
	slices.Sort(assRess)
	for _, res := range assRess {
		res := res
		btn := eeui.NewImageButton(ui.imageHandler.createThumbnailOverlay(ItemTypeAssembler, res), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ItemTypeAssembler, res)
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

	//Finalizers
	var finBtns []eeui.Widget
	for _, res := range ui.config.Resources {
		res := res
		btn := eeui.NewImageButton(ui.imageHandler.createThumbnailOverlay(ItemTypeFinalizer, res), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ItemTypeFinalizer, res)
		})
		finBtns = append(finBtns, btn)
	}
	finLayout := eeui.NewHBoxLayout(
		eeui.BoxLayoutStyles{
			Padding: 4,
			Gap:     4,
			SizeHint: eeui.SizeHint{
				MaxHeight: 48,
			},
		},
		finBtns...,
	)

	//Misc
	var miscBtns []eeui.Widget
	{
		btn := eeui.NewImageButton(ui.assets.Trash(), 48, 48, evts)
		btn.OnClick(func() {
			selectItem(ItemTypeTrash, minifac.NoResource)
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
		finLayout,
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
			obj, err := ui.CreateObject(ui.selectedItem, ui.selectedResource)
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
	assets           *assets.Assets
	config           *minifac.Config
	universe         *minifac.Universe
	backgroundImg    *ebiten.Image
	imageHandler     *ImageHandler
	ticker           *time.Ticker
	running          bool
	menu             *eeui.Form
	selectedItem     ItemType
	selectedResource minifac.Resource
}

func (ui *UI) createBackground() {
	img := ebiten.NewImage(ui.dx, ui.dy)
	sz := ui.universe.Size()
	for ux := 0; ux < sz.DX; ux++ {
		x := ui.scaleX * float64(ux)
		vector.StrokeLine(img, float32(x), 0, float32(x), float32(ui.dy), 1, color.RGBA{0, 0, 255, 255}, true)
	}
	for uy := 0; uy < sz.DY; uy++ {
		y := ui.scaleY * float64(uy)
		vector.StrokeLine(img, 0, float32(y), float32(ui.dx), float32(y), 1, color.RGBA{0, 0, 255, 255}, true)
	}
	ui.backgroundImg = img
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
	dx := outsideWidth - MenuWidth
	dy := outsideHeight
	if dx == ui.dx && dy == ui.dy {
		return outsideWidth, outsideHeight
	}

	ui.backgroundImg = nil
	ui.menu.Resize(outsideWidth-MenuWidth, 0, MenuWidth, outsideHeight)
	ui.dx = dx
	ui.dy = dy

	sz := ui.universe.Size()
	ui.scaleX, ui.scaleY = float64(ui.dx)/float64(sz.DX), float64(ui.dy)/float64(sz.DY)
	return outsideWidth, outsideHeight
}

func (ui *UI) Draw(screen *ebiten.Image) {
	if ui.backgroundImg == nil {
		ui.createBackground()
	}
	screen.DrawImage(ui.backgroundImg, &ebiten.DrawImageOptions{})

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
