package main

import (
	"fmt"

	"github.com/eihigh/canvas"
	renderer "github.com/eihigh/canvas/renderers/ebiten"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
)

type scene int

const (
	title scene = iota
	play
	gameover
)

const (
	ps     = 50.0 // player size
	vw, vh = ps * 8, ps * 12
)

func (g *game) draw(screen *ebiten.Image) {
	r := renderer.New(screen)
	c := canvas.NewContext(r)

	// Calc scroll
	scroll := 0.0
	realY := g.y + jumpHeight*g.t
	if realY > 4 {
		scroll = realY - 4
	}

	// Draw background
	{
		op := &ebiten.DrawImageOptions{}
		img := images["assets/sky.png"]
		w, h := img.Size()
		sw, sh := screen.Size()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(1.2, 1.2)
		op.GeoM.Translate(float64(sw)/2, float64(sh)/2)
		screen.DrawImage(img, op)
	}

	// Draw pillars
	{
		hideLeft := g.left
		switch g.state {
		case wall:
			hideLeft = !hideLeft
		}
		m := int(g.y/brickSize) - 10
		if m < 0 {
			m = 0
		}
		n := m + 20
		if n > len(g.leftPillar) {
			n = len(g.leftPillar)
		}

		for i := m; i < n; i++ {
			b := g.leftPillar[i]
			img := images[b.imgName()]
			y := float64(i)*brickSize - scroll
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, vh-y*ps-ps*2)
			if hideLeft {
				op.ColorM.ChangeHSV(1, 1, 0.3)
			}
			screen.DrawImage(img, op)
		}

		for i := m; i < n; i++ {
			b := g.rightPillar[i]
			img := images[b.imgName()]
			y := float64(i)*brickSize - scroll
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(vw-brickSize*ps, vh-y*ps-ps*2)
			if !hideLeft {
				op.ColorM.ChangeHSV(1, 1, 0.3)
			}
			screen.DrawImage(img, op)
		}
	}

	// Draw jump result
	if len(g.jumpLog) > 0 {
		tail := g.jumpLog[len(g.jumpLog)-1]
		if (tick - tail.reachAt) < 60 {
			s := "OK!"
			if !tail.ok {
				s = "Miss..."
			}
			face := font.Face(30, canvas.Black, canvas.FontRegular, canvas.FontNormal)
			face.FauxBold = 1.0 / 30
			text := canvas.NewTextBox(face, s, vw, vh/2, canvas.Center, canvas.Center, 0, 0)
			c.DrawText(0, 0, text)

			face = font.Face(30, canvas.White, canvas.FontRegular, canvas.FontNormal)
			text = canvas.NewTextBox(face, s, vw, vh/2, canvas.Center, canvas.Center, 0, 0)
			c.DrawText(0, 0, text)
		}
	}

	// Draw player
	if g.scene != gameover {
		op := &ebiten.DrawImageOptions{}
		img := images["assets/man.png"]
		size, h := img.Size()
		if h > size {
			size = h
		}
		sx, sy := ps/float64(size), ps/float64(size)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(0, -float64(h)*sy)

		// Calc y
		y := g.y - scroll
		if g.state == air {
			y += jumpHeight * ease.OutQuad(g.t)
		}
		// Calc x
		x := 0.0
		tx := g.t
		if !g.left {
			x = 3.0
			tx = -tx
		}
		if g.state == air {
			x += 3.0 * tx
		}

		op.GeoM.Translate(ps*2+x*ps, vh-y*ps)
		screen.DrawImage(img, op)
	}

	// Draw fire
	{
		op := &ebiten.DrawImageOptions{}
		fireY := g.fireY - scroll
		img := images["assets/bg_fire.jpeg"]
		w, h := img.Size()
		sw, sh := screen.Size()
		sx := float64(sw) / float64(w)
		sy := float64(sh) / float64(h)
		op.GeoM.Scale(sx, sy)
		y := vh/4 - fireY*ps
		op.GeoM.Translate(0, y)
		op.CompositeMode = ebiten.CompositeModeLighter
		screen.DrawImage(img, op)
	}

	// Draw title message
	if g.scene == title {
		face := font.Face(30, canvas.Black, canvas.FontRegular, canvas.FontNormal)
		face.FauxBold = 2.0 / 30
		text := canvas.NewTextBox(face, "Climbing Magnet Man", vw, vh/2, canvas.Center, canvas.Center, 0, 0)
		c.DrawText(0, 0, text)

		face = font.Face(30, canvas.White, canvas.FontRegular, canvas.FontNormal)
		text = canvas.NewTextBox(face, "Climbing Magnet Man", vw, vh/2, canvas.Center, canvas.Center, 0, 0)
		c.DrawText(0, 0, text)

		face = font.Face(24, canvas.White, canvas.FontRegular, canvas.FontNormal)
		text = canvas.NewTextBox(face, "Space key / Click / Touch\nto jump.", vw, vh, canvas.Center, canvas.Center, 0, 0)
		c.DrawText(0, 0, text)
	} else {
		// Draw gameover message
		if g.scene == gameover {
			face := font.Face(30, canvas.Black, canvas.FontRegular, canvas.FontNormal)
			face.FauxBold = 2.0 / 30
			text := canvas.NewTextBox(face, "GAME OVER", vw, vh/2, canvas.Center, canvas.Center, 0, 0)
			c.DrawText(0, 0, text)

			face = font.Face(30, canvas.White, canvas.FontRegular, canvas.FontNormal)
			text = canvas.NewTextBox(face, "GAME OVER", vw, vh/2, canvas.Center, canvas.Center, 0, 0)
			c.DrawText(0, 0, text)
		}

		// Draw score
		s := fmt.Sprintf("%dm", int(g.y*2))
		face := font.Face(28, canvas.Black, canvas.FontRegular, canvas.FontNormal)
		face.FauxBold = 1.0 / 30
		text := canvas.NewTextBox(face, s, vw, vh/4, canvas.Center, canvas.Center, 0, 0)
		c.DrawText(0, vh/4, text)

		face = font.Face(28, canvas.White, canvas.FontRegular, canvas.FontNormal)
		text = canvas.NewTextBox(face, s, vw, vh/4, canvas.Center, canvas.Center, 0, 0)
		c.DrawText(0, vh/4, text)
	}

	// Draw debug message
	// dmsg := fmt.Sprint(g.fireVel(), g.fireY)
	// ebitenutil.DebugPrint(screen, dmsg)
}

func (b *brick) imgName() string {
	if b.pole() == N {
		return "gen/N"
	}
	return "gen/S"
}
