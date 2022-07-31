package raycast

import (
	"errors"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const moveAmount = 0.002
const rotateAmount = 0.0005

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

	playerPos       point
	playerDir       point
	playerStrafeDir point
	plane           point
	oldMousePos     int
}

func NewWorld(width, height int) *World {
	w := &World{
		width:  width,
		height: height,
		tiles:  make([][]*tile, width*height),
		playerDir: point{
			x: 0,
			y: -1,
		},
		plane: point{
			x: 0.66,
			y: 0,
		},
		playerStrafeDir: point{
			x: 1,
			y: 0,
		},
		playerPos: point{
			x: 3,
			y: 3,
		},
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
		movePlayer(w, delta, point{
			x: w.playerDir.x,
			y: w.playerDir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		movePlayer(w, delta, point{
			x: -w.playerDir.x,
			y: -w.playerDir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movePlayer(w, delta, point{
			x: -w.playerStrafeDir.x,
			y: -w.playerStrafeDir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		movePlayer(w, delta, point{
			x: w.playerStrafeDir.x,
			y: w.playerStrafeDir.y,
		})
	}

	// mouse look
	mx, _ := ebiten.CursorPosition()
	mouseMove := w.oldMousePos - mx
	rotation := rotateAmount * delta * float64(mouseMove)

	oldDirX := w.playerDir.x
	w.playerDir.x = w.playerDir.x*math.Cos(-rotation) - w.playerDir.y*math.Sin(-rotation)
	w.playerDir.y = oldDirX*math.Sin(-rotation) + w.playerDir.y*math.Cos(-rotation)
	oldPlaneX := w.plane.x
	w.plane.x = w.plane.x*math.Cos(-rotation) - w.plane.y*math.Sin(-rotation)
	w.plane.y = oldPlaneX*math.Sin(-rotation) + w.plane.y*math.Cos(-rotation)
	oldStrafeX := w.playerStrafeDir.x
	w.playerStrafeDir.x = w.playerStrafeDir.x*math.Cos(-rotation) - w.playerStrafeDir.y*math.Sin(-rotation)
	w.playerStrafeDir.y = oldStrafeX*math.Sin(-rotation) + w.playerStrafeDir.y*math.Cos(-rotation)
	w.oldMousePos = mx

	// syscalls
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("normal escape termination")
	}
	return nil
}

func movePlayer(w *World, delta float64, dir point) {
	movex := dir.x * moveAmount * delta
	movey := dir.y * moveAmount * delta

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

}
