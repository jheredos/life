package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	ScreenWidth  = 600
	ScreenHeight = 400
)

type Game struct {
	world  *World
	pixels []byte
}

var ticks int = 0

func (g *Game) Update() error {
	ticks = (ticks + 1) % 2
	if ticks == 0 {
		g.world.Update()
	}
	// g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, ScreenWidth*ScreenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
