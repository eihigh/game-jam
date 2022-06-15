package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	_ "image/png"
)

const (
	tau = math.Pi * 2
)

var (
	vw, vh = 400, 400
	tick   = 0
	frame  = 0
)

type game struct{}

func newGame() (*game, error) {
	g := &game{}

	// Load Assets.
	for _, name := range []string{
		"assets/jisyaku_bou.png",
	} {
		if err := loadImage(name); err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g *game) Update() error {
	tick++
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	frame++

	op := &ebiten.DrawImageOptions{}
	i := images["assets/jisyaku_bou.png"]
	w, h := i.Size()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	e := float64(tick) / 30
	op.GeoM.Rotate(tau * e)
	screen.DrawImage(images["assets/jisyaku_bou.png"], op)
}

func (g *game) Layout(ow, oh int) (int, int) {
	return vw, vh
}

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() error {
	g, err := newGame()
	if err != nil {
		return err
	}
	ebiten.SetWindowTitle("Ebitengine Game Jam")
	ebiten.SetWindowSize(vw, vh)
	return ebiten.RunGame(g)
}
