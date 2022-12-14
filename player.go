package raycast

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math"
)

const checkDistance = 0.7
const maxAmmo = 30
const maxHealth = 10
const screenFlashTime = 200 // millis

type player struct {
	pos              vector
	dir              vector
	strafeDir        vector
	plane            vector
	oldMousePos      int
	ammo             int
	fireRateTimer    float64
	fireRateMax      float64
	width            float64
	health           int
	souls            int
	weaponAnimation  *animation
	useWeaponTimer   float64
	showMiniMap      bool
	screenFlashColor color.RGBA
	screenFlashTimer float64
	oldHealth        int
	keys             int
}

func NewPlayer(pos vector, dir string) *player {
	p := &player{
		dir: vector{
			x: 1,
			y: 0,
		},
		plane: vector{
			x: 0,
			y: 0.5,
		},
		strafeDir: vector{
			x: 0,
			y: 1,
		},
		pos:         pos,
		fireRateMax: 400.0, // millis
		ammo:        30,
		width:       0.5,
		health:      3,
		oldHealth:   3,
		weaponAnimation: &animation{
			numFrames: 4,
			numTime:   0.1 * 1000,
			isReset:   true,
		},
		showMiniMap: false,
	}
	switch dir {
	case "north":
		p.dir = vector{
			x: 0,
			y: -1,
		}
		p.strafeDir = vector{
			x: 1,
			y: 0,
		}
		p.plane = vector{
			x: 0.5,
			y: 0,
		}
	case "south":
		p.dir = vector{
			x: 0,
			y: 1,
		}
		p.strafeDir = vector{
			x: -1,
			y: 0,
		}
		p.plane = vector{
			x: -0.5,
			y: 0,
		}
	case "east":
		p.dir = vector{
			x: 1,
			y: 0,
		}
		p.strafeDir = vector{
			x: 0,
			y: 1,
		}
		p.plane = vector{
			x: 0,
			y: 0.5,
		}
	case "west":
		p.dir = vector{
			x: -1,
			y: 0,
		}
		p.strafeDir = vector{
			x: 0,
			y: -1,
		}
		p.plane = vector{
			x: 0,
			y: -0.5,
		}
	}

	return p
}

func (r *player) Update(w *World, delta float64) error {
	if r.oldHealth > r.health {
		w.soundPlayer.PlaySound("player-hurt")
	}
	r.oldHealth = r.health
	r.weaponAnimation.Update(delta)

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
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		checkPos := vector{
			x: r.pos.x + (r.dir.x * checkDistance),
			y: r.pos.y + (r.dir.y * checkDistance),
		}
		t := w.getTileAtPoint(checkPos)
		if t.door {
			tryToOpenDoor(w, r, t)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		r.showMiniMap = !r.showMiniMap
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		fakeLightEnabled = !fakeLightEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		w.debug.passiveMode = !w.debug.passiveMode
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		screenFlashEnabled = !screenFlashEnabled
	}
	// change to pressed with fire rate
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// check for ammo
		if r.ammo > 0 && r.fireRateTimer < 0 {
			r.ammo -= 1
			r.fireRateTimer = r.fireRateMax
			posInFrontOfPlayer := addVector(r.pos, scaleVector(r.dir, 0.3))
			w.ShootBullet(posInFrontOfPlayer, r.dir, bulletSpeed)
			r.weaponAnimation.Play()
			w.soundPlayer.PlaySound("crack")
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

	if r.screenFlashTimer > 0 {
		r.screenFlashTimer -= delta
	}

	// syscalls
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("normal escape termination")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if r.health < 0 {
		return errors.New("player died")
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

func (r *player) TakeDamage(amount int) {
	r.health -= 1
	r.screenFlashTimer = screenFlashTime
	r.screenFlashColor = hurtScreenFlashColor
}

func tryToOpenDoor(w *World, r *player, t *tile) {
	if t.block {
		if t.locked {
			if r.keys <= 0 {
				w.soundPlayer.PlaySound("thud")
				return
			}
			r.keys -= 1
			//w.soundPlayer.PlaySound("door") // play click sound
		}
		t.block = false
		w.soundPlayer.PlaySound("door")
	} else {
		playerT := w.getTileAtPoint(r.pos)
		if playerT != t {
			t.block = true
			w.soundPlayer.PlaySound("door")
		}
	}
}
