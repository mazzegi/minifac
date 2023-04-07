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
	ImageTypeFinalizer      ImageType = "finalizer.png"
	ImageTypeConveyor_east  ImageType = "conveyor_east.png"
	ImageTypeConveyor_north ImageType = "conveyor_north.png"
	ImageTypeConveyor_south ImageType = "conveyor_south.png"
	ImageTypeConveyor_west  ImageType = "conveyor_west.png"
	ImageTypeWood           ImageType = "wood.png"
	ImageTypeCoal           ImageType = "coal.png"
	ImageTypeStone          ImageType = "stone.png"
	ImageTypeIronOre        ImageType = "ironore.png"
	ImageTypeIron           ImageType = "iron.png"
	ImageTypeSteel          ImageType = "steel.png"
)

var allImageTypes = []ImageType{
	ImageTypeProducer,
	ImageTypeAssembler,
	ImageTypeTrash,
	ImageTypeFinalizer,
	ImageTypeConveyor_east,
	ImageTypeConveyor_north,
	ImageTypeConveyor_south,
	ImageTypeConveyor_west,
	ImageTypeWood,
	ImageTypeCoal,
	ImageTypeStone,
	ImageTypeIronOre,
	ImageTypeIron,
	ImageTypeSteel,
}

func resourceImageType(res minifac.Resource) ImageType {
	switch res {
	case minifac.Wood:
		return ImageTypeWood
	case minifac.Stone:
		return ImageTypeStone
	case minifac.Coal:
		return ImageTypeCoal
	case minifac.IronOre:
		return ImageTypeIronOre
	case minifac.Iron:
		return ImageTypeIron
	case minifac.Steel:
		return ImageTypeSteel
	default:
		return ""
	}
}

type PositionedImage struct {
	Position grid.Position
	Image    *ebiten.Image
}

func mustLoadImage(typ ImageType) *ebiten.Image {
	path := filepath.Join("assets", string(typ))
	return mustLoadImageAsset(path)
}

func NewImageHandler(u *minifac.Universe) *ImageHandler {
	ih := &ImageHandler{
		universe:          u,
		images:            make(map[ImageType]*ebiten.Image),
		overlays:          make(map[imageOverlay]*ebiten.Image),
		thumbnailOverlays: make(map[imageOverlay]*ebiten.Image),
	}
	for _, it := range allImageTypes {
		img := mustLoadImage(it)
		ih.images[it] = img
	}
	return ih
}

type imageOverlay struct {
	imageType   ImageType
	overlayType ImageType
}

type ImageHandler struct {
	universe          *minifac.Universe
	images            map[ImageType]*ebiten.Image
	overlays          map[imageOverlay]*ebiten.Image
	thumbnailOverlays map[imageOverlay]*ebiten.Image
}

func (h *ImageHandler) Images() []*PositionedImage {
	//TODO: cache images
	imgs := []*PositionedImage{}
	for _, gobj := range h.universe.AllObjects() {
		switch obj := gobj.Value.(type) {
		case *minifac.IncarnationProducer:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ImageTypeProducer, resourceImageType(obj.Resource())),
			})
		case *minifac.Trashbin:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.images[ImageTypeTrash],
			})
		case *minifac.Finalizer:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ImageTypeFinalizer, resourceImageType(obj.Resource())),
			})
		case *minifac.Assembler:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ImageTypeAssembler, resourceImageType(obj.Resource())),
			})
		case *minifac.Conveyor:
			var convType ImageType
			switch obj.Dir() {
			case grid.East:
				convType = ImageTypeConveyor_east
			case grid.South:
				convType = ImageTypeConveyor_south
			case grid.West:
				convType = ImageTypeConveyor_west
			case grid.North:
				convType = ImageTypeConveyor_north
			default:
				panic(fmt.Errorf("unknown direction %v", obj.Dir()))
			}
			img := h.createOverlay(convType, resourceImageType(obj.Resource()))
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

func (h *ImageHandler) createOverlay(baseType ImageType, overlayType ImageType) *ebiten.Image {
	if img, ok := h.overlays[imageOverlay{baseType, overlayType}]; ok {
		return img
	}

	base := h.images[baseType]
	overlay, ok := h.images[overlayType]
	if !ok {
		return base
	}
	img := ebiten.NewImageFromImage(base)
	baseBounds := base.Bounds()
	overlayBounds := overlay.Bounds()
	x := (baseBounds.Dx() - overlayBounds.Dx()) / 2
	y := (baseBounds.Dy() - overlayBounds.Dy()) / 2
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	img.DrawImage(overlay, opts)
	h.overlays[imageOverlay{baseType, overlayType}] = img
	return img
}

func (h *ImageHandler) createThumbnailOverlay(baseType ImageType, overlayType ImageType) *ebiten.Image {
	if img, ok := h.thumbnailOverlays[imageOverlay{baseType, overlayType}]; ok {
		return img
	}
	base := h.images[baseType]
	overlay := h.images[overlayType]

	img := ebiten.NewImageFromImage(base)
	baseBounds := base.Bounds()

	//overlay should be in the top half of base
	overlayHeight := baseBounds.Dy() / 2
	scaleY := float64(overlayHeight) / float64(overlay.Bounds().Dy())
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleY, scaleY)
	opts.GeoM.Translate(4, 4)
	img.DrawImage(overlay, opts)
	h.thumbnailOverlays[imageOverlay{baseType, overlayType}] = img
	return img
}
