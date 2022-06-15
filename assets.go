package main

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// embed機能 - Goには実行ファイルにさまざまなデータを埋め込めるembed機能があります。
// ここでは assets ディレクトリの内容を FS 変数に埋め込んでいます。
// embed.FS 型の変数の宣言の上に //go:embed （goの前にスペースがない）を書いて
// コンパイラに指示を出すことで使えます。
// 他にもファイルとしてではなく文字列として埋め込む機能などもありますが、割愛します。

//go:embed assets
var FS embed.FS

// ファイルパスをキー、画像データを値としたmapに、ゲーム中で使う画像すべてを格納します。

var (
	images = map[string]*ebiten.Image{}
)

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
