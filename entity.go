package raycast

import (
	"math"
)

type entity struct {
	sprites       []*sprite
	pos           vector
	width         float64
	dir           vector
	speed         float64
	health        int
	state         EntityState
	currentSprite int
}

type EntityState string

const (
	DeadEntityState    EntityState = "dead"
	NothingEntityState EntityState = "nothing"
)

const entitySpeed = 0.002

func NewEntity(img string, pos vector) *entity {
	return &entity{
		sprites: []*sprite{
			{
				image: img,
				pos:   vector{},
			},
		},
		pos:    pos,
		dir:    vector{},
		speed:  entitySpeed,
		health: 2,
		state:  NothingEntityState,
		width:  (1.0 / TextureWidth) * 20.0,
	}
}

func (r *entity) Update(delta float64) {
	if r.state == DeadEntityState {
		return
	}
	if r.health < 0 {
		r.state = DeadEntityState
		// play dying animation
		// spawn pickup
	}
	r.pos.x = r.pos.x + (r.dir.x * delta * r.speed)
	r.pos.y = r.pos.y + (r.dir.y * delta * r.speed)
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
	r.sprites[r.currentSprite].animation.currentFrame = 0
	r.sprites[r.currentSprite].animation.currentTime = 0
}

func collides(e1, e2 *entity) bool {
	if e1.state == DeadEntityState || e2.state == DeadEntityState {
		return false
	}
	withinX := math.Abs(e1.pos.x-e2.pos.x) < ((e1.width + e2.width) / 2)
	withinY := math.Abs(e1.pos.y-e2.pos.y) < ((e1.width + e2.width) / 2)
	return withinX && withinY
}
