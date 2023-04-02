package ui

import (
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var assetsFS embed.FS

func loadImageAsset(path string) (*ebiten.Image, error) {
	f, err := assetsFS.Open(path)
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
