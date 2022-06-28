package main

import (
	"bytes"
	"embed"
	"image"

	"github.com/eihigh/canvas"
	renderer "github.com/eihigh/canvas/renderers/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

// embed機能 - Goには実行ファイルにさまざまなデータを埋め込めるembed機能があります。
// ここでは assets ディレクトリの内容を FS 変数に埋め込んでいます。
// embed.FS 型の変数の宣言の上に //go:embed （goの前にスペースがない）を書いて
// コンパイラに指示を出すことで使えます。
// 他にもファイルとしてではなく文字列として埋め込む機能などもありますが、割愛します。

//go:embed assets
var FS embed.FS

var (
	// ファイルパスをキー、画像データを値としたmapに、ゲーム中で使う画像すべてを格納します。
	images = map[string]*ebiten.Image{}

	font *canvas.FontFamily
)

func loadAssets() error {
	// Load font
	font = canvas.NewFontFamily("times")
	if err := font.LoadFontFileFS(FS, "assets/NimbusRomNo9L-Reg.otf", canvas.FontRegular); err != nil {
		return err
	}

	// Load images
	for _, img := range []string{
		"assets/man.png",
		"assets/sky.png",
		"assets/bg_fire.jpeg",
	} {
		if err := loadImage(img); err != nil {
			return err
		}
	}
	return nil
}

func makeAssets() error {
	bw, bh := ps*brickSize, ps*brickSize

	// Make brick images
	face := font.Face(bh*0.6, canvas.White, canvas.FontRegular, canvas.FontNormal)
	N := ebiten.NewImage(int(bw), int(bh))
	r := renderer.New(N)
	c := canvas.NewContext(r)
	c.SetFillColor(canvas.Red)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(1)
	c.DrawPath(0, 0, canvas.Rectangle(bw, bh))
	text := canvas.NewTextBox(face, "N", bw, bh, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(0, 0, text)
	images["gen/N"] = N

	S := ebiten.NewImage(int(bw), int(bh))
	r = renderer.New(S)
	c = canvas.NewContext(r)
	c.SetFillColor(canvas.Blue)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(1)
	c.DrawPath(0, 0, canvas.Rectangle(bw, bh))
	text = canvas.NewTextBox(face, "S", bw, bh, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(0, 0, text)
	images["gen/S"] = S

	return nil
}

func loadImage(name string) error {
	// embed.FSのメソッド ReadFile は、ファイルの内容を丸ごと []byte として読み込みます。
	b, err := FS.ReadFile(name)
	if err != nil {
		return err
	}
	// データを画像にします。image.Decode 関数は io.Reader を要求するので、
	// []byte を io.Reader として扱うために bytes.NewReader 関数を使います。
	i, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}
	// 画像をEbitenの画像データに変換します。
	// （内部的には、メインメモリからグラフィックメモリへデータがコピーされている）
	e := ebiten.NewImageFromImage(i)
	images[name] = e
	return nil
}
