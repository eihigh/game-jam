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

	crtShader *ebiten.Shader
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

func (g *game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	if crtShader == nil {
		var err error
		crtShader, err = ebiten.NewShader([]byte(crtText))
		if err != nil {
			panic(err)
		}
	}
	b := offscreen.Bounds()
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = offscreen
	opts.GeoM = geoM
	screen.DrawRectShader(b.Dx(), b.Dy(), crtShader, opts)
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

const crtText = `
//go:build ignore

//kage:unit pixels

// Based on https://github.com/XorDev/1PassBlur
// https://www.shadertoy.com/view/DtscRf

package main

const radius = 16.0
const samples = 32.0
const base = 0.5
const glow = 1.1

func hash2(p vec2) vec2 {
  return normalize(
    fract(cos(p * mat2(195,174,286,183)) * 742) - 0.5)
}

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
  blur := vec4(0)
  weights := 0.0
  scale := radius / sqrt(samples)
  offset := hash2(srcPos) * scale
  rot := mat2(-0.7373688, -0.6754904, 0.6754904, -0.7373688)
  for i := 0.0; i < samples; i += 1 {
    // rotate by golden angle
    offset *= rot
    dist := sqrt(i)
    pos := srcPos + offset*dist
    color := imageSrc0At(pos)

    weight := 1.0 / (1+dist)
    blur += color * weight
    weights += weight
  }
  blur /= weights
  clr := mix(blur*glow, imageSrc0At(srcPos), base)

  rgb := clr.rgb
  rgb = clamp(mix(rgb, rgb*rgb, 0.4), 0, 1)

  // vignette
  origin, size := imageSrcRegionOnTexture()
  uv := (srcPos - origin) / size
  vig := 40*uv.x*uv.y*(1-uv.x)*(1-uv.y)
  rgb *= vec3(pow(vig, 0.3))
  rgb *= vec3(0.95, 1.05, 0.95)

  rgb *= 1.0 - mod(dstPos.y, 2)*0.2

  rgb *= 1.4
  return vec4(rgb, clr.a)
}
`
