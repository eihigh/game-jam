package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	_ "image/jpeg"
	_ "image/png"
)

const (
	// 個人的な趣味で τ を使います。
	tau = math.Pi * 2
)

var (
	// Updateが呼ばれることを、Ebitenでは tick と呼びます。
	tick = 0
	// Draw が呼ばれることを、Ebitenでは frame と呼びます。
	frame = 0
)

// ゲームを初期化します。
func newGame() (*game, error) {
	g := &game{}
	// Load Assets.
	if err := loadAssets(); err != nil {
		return nil, err
	}
	g.restart()
	return g, nil
}

// Update関数は、画面のリフレッシュレートに関わらず
// 常に毎秒60回呼ばれます（既定値）。
// 描画ではなく更新処理を行うことが推奨されます。
func (g *game) Update() error {
	if tick == 0 {
		// makeAssets is only available after RunGame
		if err := makeAssets(); err != nil {
			return err
		}
	}
	tick++

	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0

	i := &input{
		jump: jump,
	}

	return g.update(i)
}

// Draw関数は、画面のリフレッシュレートと同期して呼ばれます（既定値）。
// 描画処理のみを行うことが推奨されます。ここで状態の変更を行うといろいろ事故ります。
func (g *game) Draw(screen *ebiten.Image) {
	frame++
	g.draw(screen)
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
	ebiten.SetWindowTitle("Climbing Magnet Man")
	// ウィンドウサイズを決定します。
	ebiten.SetWindowSize(vw, vh)
	// ゲームスタート！
	return ebiten.RunGame(g)
}
