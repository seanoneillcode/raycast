package raycast

const gravity = 0.0009

type particleType string

const (
	smokeParticleType particleType = "smokeParticle"
)

type particle struct {
	sprite    *sprite
	pos       vector
	acc       vector
	heightAcc float64
	height    float64
	speed     float64
	ttl       float64
}

func NewParticle(pos vector, acc vector, height float64, heightAcc float64, speed float64, ttl float64, sprite *sprite) *particle {
	return &particle{
		pos:       pos,
		sprite:    sprite,
		height:    height,
		acc:       acc,
		heightAcc: heightAcc,
		speed:     speed,
		ttl:       ttl,
	}
}

func (r *particle) Update(delta float64, w *World) {
	r.ttl = r.ttl - delta
	if r.ttl < 0 {
		return
	}

	r.pos.x = r.pos.x + (r.acc.x * delta * r.speed)
	r.pos.y = r.pos.y + (r.acc.y * delta * r.speed)
	r.height = r.height + r.heightAcc

	// dampen
	r.heightAcc = r.heightAcc + gravity*delta
	if r.heightAcc > gravity {
		r.heightAcc = gravity
	}
	r.acc.x = r.acc.x * physicsDampening
	r.acc.y = r.acc.y * physicsDampening

	r.sprite.pos.x = r.pos.x
	r.sprite.pos.y = r.pos.y
	r.sprite.height = r.height
	if r.sprite.animation != nil {
		r.sprite.animation.Update(delta)
	}
}
