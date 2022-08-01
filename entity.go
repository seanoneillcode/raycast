package raycast

type sprite struct {
	image    string
	pos      vector
	distance float64
	height   float64
}

type entity struct {
	sprite *sprite
	pos    vector
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

func NewEntity(img string, pos vector) *entity {
	return &entity{
		sprite: &sprite{
			image: img,
			pos:   vector{},
		},
		pos:    pos,
		dir:    vector{},
		speed:  0.002,
		health: 2,
		state:  NothingEntityState,
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

func (r *bullet) Update(w *World, delta float64) {
	r.entity.Update(delta)
	t := w.getTileAtPoint(r.entity.pos)
	if t.block {
		r.entity.state = DeadEntityState
		r.entity.undoLastMove(delta)
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
