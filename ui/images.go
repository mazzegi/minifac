package ui

import (
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/grid"
)

type ImageType string

const (
	ImageTypeProducer       ImageType = "producer.png"
	ImageTypeAssembler      ImageType = "assembler.png"
	ImageTypeTrash          ImageType = "trash.png"
	ImageTypeConveyor_east  ImageType = "conveyor_east.png"
	ImageTypeConveyor_north ImageType = "conveyor_north.png"
	ImageTypeConveyor_south ImageType = "conveyor_south.png"
	ImageTypeConveyor_west  ImageType = "conveyor_west.png"
	ImageTypeCoal           ImageType = "coal.png"
	ImageTypeIron           ImageType = "iron.png"
	ImageTypeSteel          ImageType = "steel.png"
)

var allImageTypes = []ImageType{
	ImageTypeProducer,
	ImageTypeAssembler,
	ImageTypeTrash,
	ImageTypeConveyor_east,
	ImageTypeConveyor_north,
	ImageTypeConveyor_south,
	ImageTypeConveyor_west,
	ImageTypeCoal,
	ImageTypeIron,
	ImageTypeSteel,
}

type PositionedImage struct {
	Position grid.Position
	Image    *ebiten.Image
}

func NewImageHandler(u *minifac.Universe) *ImageHandler {
	ih := &ImageHandler{
		universe: u,
		images:   make(map[ImageType]*ebiten.Image),
	}
	for _, it := range allImageTypes {
		path := filepath.Join("assets", string(it))
		img, err := loadImageAsset(path)
		if err != nil {
			panic(fmt.Errorf("loading image %q", path))
		}
		ih.images[it] = img
	}
	return ih
}

type ImageHandler struct {
	universe *minifac.Universe
	images   map[ImageType]*ebiten.Image
}

func (h *ImageHandler) Images() []*PositionedImage {
	imgs := []*PositionedImage{}

	for _, gobj := range h.universe.AllObjects() {
		switch obj := gobj.Value.(type) {
		case *minifac.IncarnationProducer:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.images[ImageTypeProducer],
			})
		case *minifac.Trashbin:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.images[ImageTypeTrash],
			})
		case *minifac.Assembler:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.images[ImageTypeAssembler],
			})
		case *minifac.Conveyor:
			var img *ebiten.Image
			switch obj.Dir() {
			case grid.East:
				img = h.images[ImageTypeConveyor_east]
			case grid.South:
				img = h.images[ImageTypeConveyor_south]
			case grid.West:
				img = h.images[ImageTypeConveyor_west]
			case grid.North:
				img = h.images[ImageTypeConveyor_north]
			default:
				panic(fmt.Errorf("unknown direction %v", obj.Dir()))
			}
			res := obj.Resource()
			switch res {
			case minifac.Coal:
				img = createOverlayImage(img, h.images[ImageTypeCoal])
			case minifac.Iron:
				img = createOverlayImage(img, h.images[ImageTypeIron])
			case minifac.Steel:
				img = createOverlayImage(img, h.images[ImageTypeSteel])
			}

			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    img,
			})
		default:
			panic(fmt.Errorf("unknown object type %T", obj))
		}
	}

	return imgs
}

func createOverlayImage(base *ebiten.Image, overlay *ebiten.Image) *ebiten.Image {
	img := ebiten.NewImageFromImage(base)
	baseBounds := base.Bounds()
	overlayBounds := overlay.Bounds()
	x := (baseBounds.Dx() - overlayBounds.Dx()) / 2
	y := (baseBounds.Dy() - overlayBounds.Dy()) / 2
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	img.DrawImage(overlay, opts)
	return img
}
