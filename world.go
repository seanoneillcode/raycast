package raycast

import (
	"errors"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const moveAmount = 0.004
const rotateAmount = 0.001

type point struct {
	x float64
	y float64
}

type pos struct {
	x int
	y int
}

type tile struct {
	block bool
}

type World struct {
	width  int
	height int
	tiles  [][]*tile

	playerPos   point
	playerDir   float64
	oldMousePos int
}

func NewWorld(width, height int) *World {
	w := &World{
		width:  width,
		height: height,
		tiles:  make([][]*tile, width*height),
	}
	for x := 0; x < w.width; x++ {
		w.tiles[x] = make([]*tile, width*height)
		for y := 0; y < w.height; y++ {
			w.tiles[x][y] = &tile{}
		}
	}
	initWorld(w)
	return w
}

func (w *World) Update(delta float64) error {

	// handle input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		angle := w.playerDir * math.Pi / 3

		movex := math.Cos(angle) * delta * moveAmount
		movey := math.Sin(angle) * delta * moveAmount

		newposx := w.playerPos.x + (movex * PlayerWidth)
		newposy := w.playerPos.y

		tilex := w.getTile(int(newposx), int(newposy))
		if tilex == nil || !tilex.block {
			w.playerPos.x += movex
		}

		newposx = w.playerPos.x
		newposy = w.playerPos.y + (movey * PlayerWidth)

		tiley := w.getTile(int(newposx), int(newposy))
		if tiley == nil || !tiley.block {
			w.playerPos.y += movey
		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		angle := w.playerDir * math.Pi / 3

		movex := math.Cos(angle) * delta * moveAmount
		movey := math.Sin(angle) * delta * moveAmount

		newposx := w.playerPos.x - (movex * PlayerWidth)
		newposy := w.playerPos.y

		tilex := w.getTile(int(newposx), int(newposy))
		if tilex == nil || !tilex.block {
			w.playerPos.x -= movex
		}

		newposx = w.playerPos.x
		newposy = w.playerPos.y - (movey * PlayerWidth)

		tiley := w.getTile(int(newposx), int(newposy))
		if tiley == nil || !tiley.block {
			w.playerPos.y -= movey
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		angle := w.playerDir * math.Pi / 3

		movex := math.Sin(angle) * delta * moveAmount
		movey := -math.Cos(angle) * delta * moveAmount

		newposx := w.playerPos.x + (movex * PlayerWidth)
		newposy := w.playerPos.y

		tilex := w.getTile(int(newposx), int(newposy))
		if tilex == nil || !tilex.block {
			w.playerPos.x += movex
		}

		newposx = w.playerPos.x
		newposy = w.playerPos.y + (movey * PlayerWidth)

		tiley := w.getTile(int(newposx), int(newposy))
		if tiley == nil || !tiley.block {
			w.playerPos.y += movey
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		angle := w.playerDir * math.Pi / 3

		movex := -math.Sin(angle) * delta * moveAmount
		movey := math.Cos(angle) * delta * moveAmount

		newposx := w.playerPos.x + (movex * PlayerWidth)
		newposy := w.playerPos.y

		tilex := w.getTile(int(newposx), int(newposy))
		if tilex == nil || !tilex.block {
			w.playerPos.x += movex
		}

		newposx = w.playerPos.x
		newposy = w.playerPos.y + (movey * PlayerWidth)

		tiley := w.getTile(int(newposx), int(newposy))
		if tiley == nil || !tiley.block {
			w.playerPos.y += movey
		}
	}

	// mouse look
	mx, _ := ebiten.CursorPosition()
	mouseMove := w.oldMousePos - mx
	w.playerDir += rotateAmount * delta * float64(-mouseMove)
	if w.playerDir < 0 {
		w.playerDir += math.Pi * 2
	}
	if w.playerDir > math.Pi*2 {
		w.playerDir -= math.Pi * 2
	}
	w.oldMousePos = mx

	// syscalls
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("normal escape termination")
	}
	return nil
}

func (w *World) getTileAtPoint(x, y float64) *tile {
	return w.getTile(int(x), int(y))
}

func (w *World) getTile(x, y int) *tile {
	if x < 0 || x > w.width-1 {
		return nil
	}
	if y < 0 || y > w.height-1 {
		return nil
	}
	return w.tiles[x][y]
}

func initWorld(w *World) {

	nums := [][]uint8{
		{1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1},
	}

	for ix, x := range nums {
		for iy, y := range x {
			if y == 1 {
				w.tiles[ix][iy].block = true
			}
		}
	}

	w.playerPos.x = 2
	w.playerPos.y = 4
}
