package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
)

// World represents the game state.
type World struct {
	area   []bool
	width  int
	height int
}

type Color struct {
	R byte
	G byte
	B byte
	A byte
}

var LiveColor = Color{
	R: 0x10,
	G: 0xd0,
	B: 0xa0,
	A: 0xff,
}

var DeadColor = Color{
	R: 0x00,
	G: 0x00,
	B: 0x00,
	A: 0x00,
}

// init inits world with a random state.
func (w *World) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x] = true
	}
}

// Update game state by one tick.
func (w *World) Update() {
	width := w.width
	height := w.height
	next := make([]bool, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			population := neighbourCount(w.area, width, height, x, y)
			switch {
			case population < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[y*width+x] = false
			case (population == 2 || population == 3) && w.area[y*width+x]:
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[y*width+x] = true
			case population > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[y*width+x] = false
			case population == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[y*width+x] = true
			}
		}
	}

	randomLive := 100
	for randomLive > 0 {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		if neighbourCount(w.area, width, height, x, y) > 0 {
			next[y*w.width+x] = true
		}
		randomLive--
	}

	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		next[my*width+mx] = true
	}

	w.area = next
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			pix[4*i] = LiveColor.R
			pix[4*i+1] = LiveColor.G
			pix[4*i+2] = LiveColor.B
			pix[4*i+3] = LiveColor.A
		} else {
			// pix[4*i] = DeadColor.R
			// pix[4*i+1] = DeadColor.G
			// pix[4*i+2] = DeadColor.B
			// pix[4*i+3] = DeadColor.A
			pix[4*i] = 0
			pix[4*i+1] = byte(float64(pix[4*i+1]) * .95)
			pix[4*i+2] = byte(float64(pix[4*i+2]) * .9999) // blue fades slower than green
			pix[4*i+3] = 0
		}
	}
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]bool, width*height),
		width:  width,
		height: height,
	}
	w.init(maxInitLiveCells)
	return w
}

// neighbourCount calculates the Moore neighborhood of (x, y).
func neighbourCount(area []bool, width, height, x, y int) int {
	neighbors := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			if area[y2*width+x2] {
				neighbors++
			}
		}
	}
	return neighbors
}
