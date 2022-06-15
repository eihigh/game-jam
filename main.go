package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	_ "image/png"
)

const (
	// 個人的な趣味で τ を使います。
	tau = math.Pi * 2
)

var (
	// 画面サイズです。どこでも使えるようにグローバル変数にすると便利です。
	vw, vh = 400, 400
	// Updateが呼ばれることを、Ebitenでは tick と呼びます。
	tick = 0
	// Draw が呼ばれることを、Ebitenでは frame と呼びます。
	frame = 0
)

// ebiten.Game interface を満たす型がEbitenには必要なので、この game 構造体に
// Update, Draw, Layout 関数を持たせます。
type game struct{}

// ゲームを初期化します。
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

// Update関数は、画面のリフレッシュレートに関わらず
// 常に毎秒60回呼ばれます（既定値）。
// 描画ではなく更新処理を行うことが推奨されます。
func (g *game) Update() error {
	tick++
	return nil
}

// Draw関数は、画面のリフレッシュレートと同期して呼ばれます（既定値）。
// 描画処理のみを行うことが推奨されます。ここで状態の変更を行うといろいろ事故ります。
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

// Layout関数は、ウィンドウのリサイズの挙動を決定します。とりあえず常に画面サイズを返せば無難です。
func (g *game) Layout(ow, oh int) (int, int) {
	return vw, vh
}

// すべてのGoプログラムはmain packageのmain関数から始まります。
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
	// ウィンドウタイトルを変更します。
	ebiten.SetWindowTitle("Ebitengine Game Jam")
	// ウィンドウサイズを決定します。
	ebiten.SetWindowSize(vw, vh)
	// ゲームスタート！
	return ebiten.RunGame(g)
}
