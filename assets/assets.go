package assets

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
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
	Items     map[string]string `json:"items"`
	Resources map[string]string `json:"resources"`
	Obstacles map[string]string `json:"obstacles"`
}

type Assets struct {
	Items     map[string]*ebiten.Image
	Resources map[string]*ebiten.Image
	Obstacles map[string]*ebiten.Image
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
	a.Items = make(map[string]*ebiten.Image)
	a.Resources = make(map[string]*ebiten.Image)
	a.Obstacles = make(map[string]*ebiten.Image)
	for key, path := range ac.Items {
		img, err := loadImage(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("load-image %q: %w", filepath.Join(dir, path), err)
		}
		a.Items[key] = img
	}
	for key, path := range ac.Resources {
		img, err := loadImage(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("load-image %q: %w", filepath.Join(dir, path), err)
		}
		a.Resources[key] = img
	}
	for key, path := range ac.Obstacles {
		img, err := loadImage(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("load-image %q: %w", filepath.Join(dir, path), err)
		}
		a.Obstacles[key] = img
	}
	return nil
}
