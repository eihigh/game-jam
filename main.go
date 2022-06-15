package main

import "github.com/hajimehoshi/ebiten/v2"

var (
	vw, vh = 400, 400
)

type game struct{}

func newGame() (*game, error) {
	g := &game{}

	return g, nil
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {

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
