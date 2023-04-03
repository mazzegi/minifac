package ui

import (
	"embed"
	"fmt"
	"io"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
)

//go:embed fonts
var fontsFS embed.FS

func loadFont(path string) (*truetype.Font, error) {
	f, err := fontsFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	fontBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read-font-bytes: %w", err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("parse-font: %w", err)
	}
	return font, nil
}

func mustLoadFont(path string) *truetype.Font {
	img, err := loadFont(path)
	if err != nil {
		panic(err)
	}
	return img
}
