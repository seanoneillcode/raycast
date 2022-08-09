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

type enemy struct {
	entity *entity
	// AI state
}

func (r *enemy) Update(delta float64) {
	r.entity.Update(delta)
	// AI behaviour
}

type effect struct {
	entity *entity
	timer  float64
}

func NewEffect(image string, pos vector, timing float64, numFrames int) *effect {
	e := &effect{
		entity: NewEntity(image, pos),
		timer:  float64(numFrames) * timing,
	}
	e.entity.sprite.animation = &animation{
		numFrames: numFrames,
		numTime:   timing,
		autoplay:  true,
	}
	return e
}

func (r *effect) Update(delta float64) {
	r.entity.Update(delta)
	if r.timer > 0 {
		r.timer -= delta
		if r.timer <= 0 {
			r.entity.state = DeadEntityState
		}
	}
}

type portal struct {
	entity *entity
}

func NewPortal(pos vector) *portal {
	timing := 0.2 * 1000
	p := &portal{
		entity: NewEntity("portal", pos),
	}
	p.entity.sprite.animation = &animation{
		numFrames: 4,
		numTime:   timing,
		autoplay:  true,
	}
	return p
}

func (r *portal) Update(w *World, delta float64) {
	r.entity.Update(delta)
	withinX := math.Abs(w.player.pos.x-r.entity.pos.x) < ((w.player.width + r.entity.width) / 2)
	withinY := math.Abs(w.player.pos.y-r.entity.pos.y) < ((w.player.width + r.entity.width) / 2)
	if withinX && withinY {
		fmt.Printf("entity player collided with portal\n ")
		fmt.Printf("level won!!")
		panic("player won")
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
