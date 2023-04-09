package assets

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazzegi/minifac/grid"
	"golang.org/x/exp/maps"
)

const assetsConfigFile = "minifac.assets.json"

func Load(dir string) (*Assets, error) {
	adir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("filepath.abs %q: %w", dir, err)
	}
	f, err := os.Open(filepath.Join(adir, assetsConfigFile))
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filepath.Join(adir, assetsConfigFile), err)
	}
	defer f.Close()

	var ac assetsConfig
	err = json.NewDecoder(f).Decode(&ac)
	if err != nil {
		return nil, fmt.Errorf("json.decode: %w", err)
	}
	a := &Assets{}
	err = a.load(adir, ac)
	if err != nil {
		return nil, fmt.Errorf("assets.load: %w", err)
	}
	return a, nil
}

type assetsConfig struct {
	Items struct {
		Producer  string `json:"producer"`
		Conveyor  string `json:"conveyor"`
		Assembler string `json:"assembler"`
		Finalizer string `json:"finalizer"`
		Trash     string `json:"trash"`
	} `json:"items"`
	Resources map[string]string `json:"resources"`
	Obstacles map[string]string `json:"obstacles"`
}

type Assets struct {
	items struct {
		producer  *ebiten.Image
		conveyor  *ebiten.Image
		conveyors map[grid.Direction]*ebiten.Image
		assembler *ebiten.Image
		finalizer *ebiten.Image
		trash     *ebiten.Image
	}
	resources map[string]*ebiten.Image
	obstacles map[string]*ebiten.Image
}

func (a *Assets) Conveyor(dir grid.Direction) *ebiten.Image {
	return a.items.conveyors[dir]
}

func (a *Assets) Producer() *ebiten.Image {
	return a.items.producer
}

func (a *Assets) Assembler() *ebiten.Image {
	return a.items.assembler
}

func (a *Assets) Finalizer() *ebiten.Image {
	return a.items.finalizer
}

func (a *Assets) Trash() *ebiten.Image {
	return a.items.trash
}

func (a *Assets) Resource(s string) (*ebiten.Image, bool) {
	img, ok := a.resources[s]
	return img, ok
}

func (a *Assets) Obstacle(s string) (*ebiten.Image, bool) {
	img, ok := a.obstacles[s]
	return img, ok
}

func (a *Assets) ResourceNames() []string {
	names := maps.Keys(a.resources)
	sort.Strings(names)
	return names
}

func (a *Assets) ObstacleNames() []string {
	names := maps.Keys(a.obstacles)
	sort.Strings(names)
	return names
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()
	var decodeFnc func(io.Reader) (image.Image, error)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		decodeFnc = png.Decode
	case ".jpg", ".jpeg":
		decodeFnc = jpeg.Decode
	default:
		return nil, fmt.Errorf("no decoder for extension %q", strings.ToLower(filepath.Ext(path)))
	}
	img, err := decodeFnc(f)
	if err != nil {
		return nil, fmt.Errorf("decode %q: %w", path, err)
	}
	return ebiten.NewImageFromImage(img), nil
}

func (a *Assets) load(dir string, ac assetsConfig) error {
	a.resources = make(map[string]*ebiten.Image)
	a.obstacles = make(map[string]*ebiten.Image)

	err := GroupErrors(
		func() (err error) { a.items.producer, err = loadImage(filepath.Join(dir, ac.Items.Producer)); return },
		func() (err error) { a.items.conveyor, err = loadImage(filepath.Join(dir, ac.Items.Conveyor)); return },
		func() (err error) { a.items.assembler, err = loadImage(filepath.Join(dir, ac.Items.Assembler)); return },
		func() (err error) { a.items.finalizer, err = loadImage(filepath.Join(dir, ac.Items.Finalizer)); return },
		func() (err error) { a.items.trash, err = loadImage(filepath.Join(dir, ac.Items.Trash)); return },
	)
	if err != nil {
		return err
	}

	//build directed conveyors
	a.items.conveyors = make(map[grid.Direction]*ebiten.Image)
	a.items.conveyors[grid.East] = a.items.conveyor
	confBounds := a.items.conveyor.Bounds()
	cdx, cdy := confBounds.Dx(), confBounds.Dy()
	for _, dir := range []grid.Direction{grid.South, grid.West, grid.North} {
		img := ebiten.NewImage(cdx, cdy)
		var rad float64
		switch dir {
		case grid.South:
			rad = math.Pi / 2
		case grid.West:
			rad = math.Pi
		case grid.North:
			rad = 3 * math.Pi / 2
		}
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-float64(cdx)/2, -float64(cdy)/2)
		opts.GeoM.Rotate(rad)
		opts.GeoM.Translate(float64(cdx)/2, float64(cdy)/2)
		img.DrawImage(a.items.conveyor, opts)
		a.items.conveyors[dir] = img
	}

	for key, path := range ac.Resources {
		img, err := loadImage(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("load-image %q: %w", filepath.Join(dir, path), err)
		}
		a.resources[key] = img
	}
	for key, path := range ac.Obstacles {
		img, err := loadImage(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("load-image %q: %w", filepath.Join(dir, path), err)
		}
		a.obstacles[key] = img
	}
	return nil
}
