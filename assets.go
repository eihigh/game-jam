package main

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var FS embed.FS

var (
	images = map[string]*ebiten.Image{}
)

func loadImage(name string) error {
	b, err := FS.ReadFile(name)
	if err != nil {
		return err
	}
	i, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}
	e := ebiten.NewImageFromImage(i)
	images[name] = e
	return nil
}
