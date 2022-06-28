package main

import (
	"math/rand"
)

const (
	dt          = 1.0 / 9 // Δt  for jumping
	brickSize   = 2.0
	jumpHeight  = 2.0
	G           = 1.0 / 80 // Gravitational acceleration while falling
	fallTermVel = 1.0 / 20 // Terminal velocity of falling
	penalty     = 1.0
)

type pole int

const (
	N pole = iota
	S
)

type state int

const (
	holding state = iota // holding on the wall
	air                  // not on the wall (jumping)
	wall                 // on the wall but not holding
)

type game struct {
	scene   scene
	sceneT  int
	state   state
	y       float64 // 0.0 - inf
	t       float64 // 0.0 - 1.0
	left    bool
	fallVel float64
	touchY  float64
	fireY   float64
	jumpLog []jumpLogRecord

	leftPillar, rightPillar pillar
}

type input struct {
	jump bool
}

type brick struct {
	interval int
	phase    int
}

type pillar []brick

type jumpLogRecord struct {
	reachAt int
	ok      bool
}

func (b *brick) pole() pole {
	p := (tick - b.phase) % (b.interval * 2) / b.interval
	if p == 0 {
		return N
	}
	return S
}

func (g *game) restart() {
	g.scene = title
	g.sceneT = tick
	g.state = holding
	g.y = 0
	g.t = 0
	g.left = true
	g.fallVel = 0
	g.touchY = 0
	g.fireY = -5

	g.leftPillar = nil
	g.rightPillar = nil
	g.growPillars()
}

func (g *game) update(i *input) error {
	switch g.scene {
	case title:
		// Go to play scene on the first input
		if i.jump {
			g.scene = play
			g.sceneT = tick
		}
		return nil

	case gameover:
		if i.jump && (tick-g.sceneT) > 180 {
			g.restart()
		}
		return nil
	}

	g.growPillars()

	switch g.state {
	case holding:
		// holding on the wall
		// You are ready to jump
		if i.jump {
			g.state = air
		}

	case air:
		// not on the wall (jumping)
		// Once you reach the wall, go to wall state
		g.t += dt
		if g.t >= 1.0 {
			g.t = 0
			g.left = !g.left
			g.y += jumpHeight
			g.touchY = g.y
			g.fallVel = 0
			g.state = wall
			ok := g.tryHold(false)
			r := jumpLogRecord{reachAt: tick, ok: ok}
			g.jumpLog = append(g.jumpLog, r)
		}

	case wall:
		// on the wall but not holding
		// Try to hold. If successful, go to toHold state
		if ok := g.tryHold(true); ok {
			g.state = holding
		} else {
			g.fall()
		}
	}

	// Check the deadline
	g.fireY += g.fireVel()
	if g.y < g.fireY {
		g.scene = gameover
		g.sceneT = tick
	}

	return nil
}

func (g *game) tryHold(miss bool) (ok bool) {
	if miss && (g.touchY-g.y) < penalty {
		return false
	}

	// 対象となるbrickを計算して、引き合う極同士なら成功
	i := int((g.y + 0.8) / brickSize)
	b := g.leftPillar[i]
	if !g.left {
		b = g.rightPillar[i]
	}
	p := N // left is N
	if !g.left {
		p = S // right is S
	}
	if b.pole() != p {
		g.state = holding
		return true
	}

	return false
}

func (g *game) fall() {
	g.fallVel += G
	if g.fallVel > fallTermVel {
		g.fallVel = fallTermVel
	}
	g.y -= g.fallVel
	if g.y < 0 {
		g.y = 0
		g.state = holding
	}
}

func (g *game) level(y float64) *level {
	i := 0
	for j, lvl := range levels {
		if lvl.y <= g.y {
			i = j
		}
	}
	return levels[i]
}

func (g *game) growPillars() {
	i := int(g.y / brickSize)
	l := len(g.leftPillar)
	leading := 10
	for j := l; j < i+leading; j++ {
		lvl := g.level(float64(j) * brickSize)
		interval := lvl.intervalMin + rand.Intn(lvl.intervalMax-lvl.intervalMin)
		g.leftPillar = append(g.leftPillar, brick{interval, rand.Intn(interval)})
		g.rightPillar = append(g.rightPillar, brick{interval, rand.Intn(interval)})
	}
}

func (g *game) fireVel() float64 {
	return g.level(g.y).fireVel
}

func (p pole) String() string {
	if p == N {
		return "N"
	}
	return "S"
}
