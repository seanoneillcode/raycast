package raycast

import (
	"math"
)

type entity struct {
	sprites         []*sprite
	pos             vector
	width           float64
	dir             vector
	isPhysicsEntity bool
	physics         []*vector
	speed           float64
	health          int
	state           EntityState
	currentSprite   int
	dropItem        string
}

type EntityState string

const (
	DeadEntityState    EntityState = "dead"
	NothingEntityState EntityState = "nothing"
	StunnedEntityState EntityState = "stunned"
	StoppedEntityState EntityState = "stopped"
)

const entitySpeed = 0.002
const physicsDampening = 0.9
const physicsZeroThreshold = 0.01

func NewEntity(pos vector, sprites ...*sprite) *entity {
	return &entity{
		sprites:         sprites,
		pos:             pos,
		dir:             vector{},
		speed:           entitySpeed,
		health:          1,
		state:           NothingEntityState,
		width:           (1.0 / TextureWidth) * 20.0,
		physics:         []*vector{},
		isPhysicsEntity: false,
	}
}

func (r *entity) Update(delta float64, w *World) {
	if r.state == DeadEntityState {
		return
	}
	move := vector{}
	if r.state == NothingEntityState {
		move.x = move.x + (r.dir.x * delta * r.speed)
		move.y = move.y + (r.dir.y * delta * r.speed)
	}

	if r.isPhysicsEntity && len(r.physics) > 0 {
		removeSomeVectors := false
		for _, acc := range r.physics {
			move.x = move.x + (acc.x * delta * r.speed)
			move.y = move.y + (acc.y * delta * r.speed)

			acc.x = acc.x * physicsDampening
			acc.y = acc.y * physicsDampening

			if math.Abs(acc.x) < physicsZeroThreshold {
				acc.x = 0
			}

			if math.Abs(acc.y) < physicsZeroThreshold {
				acc.y = 0
			}
			if acc.x == 0 && acc.y == 0 {
				removeSomeVectors = true
			}
		}
		if removeSomeVectors {
			var newPhysics []*vector
			for _, acc := range r.physics {
				if !(acc.x == 0 && acc.y == 0) {
					newPhysics = append(newPhysics, acc)
				}
			}
			r.physics = newPhysics
		}
	}

	if r.isPhysicsEntity {
		newposx := r.pos.x + (move.x)
		newposy := r.pos.y

		tilex := w.getTile(int(newposx), int(newposy))
		if tilex == nil || !tilex.block {
			r.pos.x += move.x
		}

		newposx = r.pos.x
		newposy = r.pos.y + (move.y)

		tiley := w.getTile(int(newposx), int(newposy))
		if tiley == nil || !tiley.block {
			r.pos.y += move.y
		}
	} else {
		r.pos.x = r.pos.x + move.x
		r.pos.y = r.pos.y + move.y
	}

	r.sprites[r.currentSprite].pos.x = r.pos.x
	r.sprites[r.currentSprite].pos.y = r.pos.y
	if r.sprites[r.currentSprite].animation != nil {
		r.sprites[r.currentSprite].animation.Update(delta)
	}
}

func (r *entity) undoLastMove(delta float64) {
	r.pos.x = r.pos.x - (r.dir.x * delta * r.speed)
	r.pos.y = r.pos.y - (r.dir.y * delta * r.speed)
	r.sprites[r.currentSprite].pos.x = r.pos.x
	r.sprites[r.currentSprite].pos.y = r.pos.y
}

func (r *entity) CurrentSprite() *sprite {
	return r.sprites[r.currentSprite]
}

func (r *entity) SetCurrentSprite(index int) {
	r.currentSprite = index
	r.sprites[index].animation.currentFrame = 0
	r.sprites[index].animation.currentTime = 0
	r.sprites[index].pos = r.pos
}

func collides(e1, e2 *entity) bool {
	if e1.state == DeadEntityState || e2.state == DeadEntityState {
		return false
	}
	withinX := math.Abs(e1.pos.x-e2.pos.x) < ((e1.width + e2.width) / 2)
	withinY := math.Abs(e1.pos.y-e2.pos.y) < ((e1.width + e2.width) / 2)
	return withinX && withinY
}

func collidesWithPlayer(player *player, e2 *entity) bool {
	if e2.state == DeadEntityState {
		return false
	}
	withinX := math.Abs(player.pos.x-e2.pos.x) < ((player.width + e2.width) / 2)
	withinY := math.Abs(player.pos.y-e2.pos.y) < ((player.width + e2.width) / 2)
	return withinX && withinY
}

func within(p1 vector, p2 vector, distance float64) bool {
	withinX := math.Abs(p1.x-p2.x) < distance
	withinY := math.Abs(p1.y-p2.y) < distance
	return withinX && withinY
}
