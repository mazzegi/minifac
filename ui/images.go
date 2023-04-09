package ui

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazzegi/minifac"
	"github.com/mazzegi/minifac/assets"
	"github.com/mazzegi/minifac/grid"
)

type PositionedImage struct {
	Position grid.Position
	Image    *ebiten.Image
}

func NewImageHandler(u *minifac.Universe, assets *assets.Assets) *ImageHandler {
	ih := &ImageHandler{
		universe:          u,
		assets:            assets,
		overlays:          make(map[imageOverlay]*ebiten.Image),
		thumbnailOverlays: make(map[imageOverlay]*ebiten.Image),
	}
	return ih
}

type imageOverlay struct {
	imageType   ItemType
	overlayType minifac.Resource
}

type ImageHandler struct {
	universe          *minifac.Universe
	assets            *assets.Assets
	overlays          map[imageOverlay]*ebiten.Image
	thumbnailOverlays map[imageOverlay]*ebiten.Image
}

func (h *ImageHandler) mustItemImage(item ItemType) *ebiten.Image {
	switch item {
	case ItemTypeProducer:
		return h.assets.Producer()
	case ItemTypeAssembler:
		return h.assets.Assembler()
	case ItemTypeTrash:
		return h.assets.Trash()
	case ItemTypeFinalizer:
		return h.assets.Finalizer()
	case ItemTypeConveyor_east:
		return h.assets.Conveyor(grid.East)
	case ItemTypeConveyor_north:
		return h.assets.Conveyor(grid.North)
	case ItemTypeConveyor_south:
		return h.assets.Conveyor(grid.South)
	case ItemTypeConveyor_west:
		return h.assets.Conveyor(grid.West)
	default:
		log.Fatalf("no image for type %q", item)
		return nil
	}
}

func (h *ImageHandler) Images() []*PositionedImage {
	imgs := []*PositionedImage{}
	for _, gobj := range h.universe.AllObjects() {
		switch obj := gobj.Value.(type) {
		case *minifac.IncarnationProducer:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ItemTypeProducer, obj.Resource()),
			})
		case *minifac.Trashbin:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.mustItemImage(ItemTypeTrash),
			})
		case *minifac.Obstacle:
			switch obj.Type() {
			default:
				if img, ok := h.assets.Obstacle(string(minifac.ObstacleWall)); ok {
					imgs = append(imgs, &PositionedImage{
						Position: gobj.Position,
						Image:    img,
					})
				}
			}
		case *minifac.Finalizer:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ItemTypeFinalizer, obj.Resource()),
			})
		case *minifac.Assembler:
			imgs = append(imgs, &PositionedImage{
				Position: gobj.Position,
				Image:    h.createThumbnailOverlay(ItemTypeAssembler, obj.Resource()),
			})
		case *minifac.Conveyor:
			var convType ItemType
			switch obj.Dir() {
			case grid.East:
				convType = ItemTypeConveyor_east
			case grid.South:
				convType = ItemTypeConveyor_south
			case grid.West:
				convType = ItemTypeConveyor_west
			case grid.North:
				convType = ItemTypeConveyor_north
			default:
				panic(fmt.Errorf("unknown direction %v", obj.Dir()))
			}
			img := h.createOverlay(convType, obj.Resource())
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

func (h *ImageHandler) createOverlay(itemType ItemType, res minifac.Resource) *ebiten.Image {
	if img, ok := h.overlays[imageOverlay{itemType, res}]; ok {
		return img
	}
	base := h.mustItemImage(itemType)
	overlay, ok := h.assets.Resource(string(res))
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
	h.overlays[imageOverlay{itemType, res}] = img
	return img
}

func (h *ImageHandler) createThumbnailOverlay(itemType ItemType, res minifac.Resource) *ebiten.Image {
	if img, ok := h.thumbnailOverlays[imageOverlay{itemType, res}]; ok {
		return img
	}
	base := h.mustItemImage(itemType)
	overlay, ok := h.assets.Resource(string(res))
	if !ok {
		return base
	}

	img := ebiten.NewImageFromImage(base)
	baseBounds := base.Bounds()

	//overlay should be in the top half of base
	overlayHeight := baseBounds.Dy() / 2
	scaleY := float64(overlayHeight) / float64(overlay.Bounds().Dy())
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleY, scaleY)
	opts.GeoM.Translate(4, 4)
	img.DrawImage(overlay, opts)
	h.thumbnailOverlays[imageOverlay{itemType, res}] = img
	return img
}
