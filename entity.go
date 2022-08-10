package raycast

import (
	"fmt"
	"math"
)

type entity struct {
	sprite *sprite
	pos    vector
	width  float64
	dir    vector
	speed  float64
	health int
	state  EntityState
}

type EntityState string

const (
	DeadEntityState    EntityState = "dead"
	NothingEntityState EntityState = "nothing"
)

const entitySpeed = 0.002

func NewEntity(img string, pos vector) *entity {
	return &entity{
		sprite: &sprite{
			image: img,
			pos:   vector{},
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
	r.sprite.pos.x = r.pos.x
	r.sprite.pos.y = r.pos.y
	if r.sprite.animation != nil {
		r.sprite.animation.Update(delta)
	}
}

func (r *entity) undoLastMove(delta float64) {
	r.pos.x = r.pos.x - (r.dir.x * delta * r.speed)
	r.pos.y = r.pos.y - (r.dir.y * delta * r.speed)
	r.sprite.pos.x = r.pos.x
	r.sprite.pos.y = r.pos.y
}

type bullet struct {
	entity *entity
}

func NewBullet(pos vector, dir vector) *bullet {
	b := &bullet{
		entity: NewEntity("bullet", pos),
	}
	b.entity.dir = dir
	b.entity.speed = bulletSpeed
	b.entity.width = (1.0 / TextureWidth) * 4.0
	return b
}

func (r *bullet) Update(w *World, delta float64) {
	r.entity.Update(delta)
	t := w.getTileAtPoint(r.entity.pos)
	if t.block {
		r.entity.state = DeadEntityState
		r.entity.undoLastMove(delta)
		w.AddEffect("bullet-hit", r.entity.pos)
	}
	for _, e := range w.enemies {
		if collides(r.entity, e.entity) {
			r.entity.state = DeadEntityState
			r.entity.undoLastMove(delta)
			w.AddEffect("bullet-hit", r.entity.pos)
			e.entity.health -= 1
			// do more here, i.e. show effects
		}
	}
}

func collides(e1, e2 *entity) bool {
	if e1.state == DeadEntityState || e2.state == DeadEntityState {
		return false
	}
	withinX := math.Abs(e1.pos.x-e2.pos.x) < ((e1.width + e2.width) / 2)
	withinY := math.Abs(e1.pos.y-e2.pos.y) < ((e1.width + e2.width) / 2)
	if withinX && withinY {
		fmt.Printf("entity %v collided with entity %v\n ", e1.sprite.image, e2.sprite.image)
	}
	return withinX && withinY
}
