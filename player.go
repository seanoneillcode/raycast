package raycast

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type player struct {
	pos           vector
	dir           vector
	strafeDir     vector
	plane         vector
	oldMousePos   int
	ammo          int
	fireRateTimer float64
	fireRateMax   float64
}

func NewPlayer(pos vector) *player {
	return &player{
		dir: vector{
			x: 0,
			y: -1,
		},
		plane: vector{
			x: 0.5,
			y: 0,
		},
		strafeDir: vector{
			x: 1,
			y: 0,
		},
		pos:         pos,
		fireRateMax: 200.0, // millis
		ammo:        100,
	}
}

func (r *player) Update(w *World, delta float64) error {
	// handle input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		r.Move(w, delta, vector{
			x: r.dir.x,
			y: r.dir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		r.Move(w, delta, vector{
			x: -r.dir.x,
			y: -r.dir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		r.Move(w, delta, vector{
			x: -r.strafeDir.x,
			y: -r.strafeDir.y,
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		r.Move(w, delta, vector{
			x: r.strafeDir.x,
			y: r.strafeDir.y,
		})
	}
	// change to pressed with fire rate
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// check for ammo
		if r.ammo > 0 && r.fireRateTimer < 0 {
			fmt.Printf("fired weapon: ammo %v; fireRateTimer: %v \n", r.ammo, r.fireRateTimer)
			r.ammo -= 1
			r.fireRateTimer = r.fireRateMax
			posInFrontOfPlayer := addVector(r.pos, r.dir)
			w.ShootBullet(posInFrontOfPlayer, r.dir)
		}
	}
	r.fireRateTimer -= delta

	// mouse look
	mx, _ := ebiten.CursorPosition()
	mouseMove := r.oldMousePos - mx
	rotation := rotateAmount * delta * float64(mouseMove)

	oldDirX := r.dir.x
	r.dir.x = r.dir.x*math.Cos(-rotation) - r.dir.y*math.Sin(-rotation)
	r.dir.y = oldDirX*math.Sin(-rotation) + r.dir.y*math.Cos(-rotation)
	oldPlaneX := r.plane.x
	r.plane.x = r.plane.x*math.Cos(-rotation) - r.plane.y*math.Sin(-rotation)
	r.plane.y = oldPlaneX*math.Sin(-rotation) + r.plane.y*math.Cos(-rotation)
	oldStrafeX := r.strafeDir.x
	r.strafeDir.x = r.strafeDir.x*math.Cos(-rotation) - r.strafeDir.y*math.Sin(-rotation)
	r.strafeDir.y = oldStrafeX*math.Sin(-rotation) + r.strafeDir.y*math.Cos(-rotation)
	r.oldMousePos = mx

	// syscalls
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("normal escape termination")
	}

	return nil
}

func (r *player) Move(w *World, delta float64, dir vector) {
	movex := dir.x * moveAmount * delta
	movey := dir.y * moveAmount * delta

	newposx := r.pos.x + (movex * PlayerWidth)
	newposy := r.pos.y

	tilex := w.getTile(int(newposx), int(newposy))
	if tilex == nil || !tilex.block {
		r.pos.x += movex
	}

	newposx = r.pos.x
	newposy = r.pos.y + (movey * PlayerWidth)

	tiley := w.getTile(int(newposx), int(newposy))
	if tiley == nil || !tiley.block {
		r.pos.y += movey
	}
}
