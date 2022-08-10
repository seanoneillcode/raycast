package raycast

type enemy struct {
	entity          *entity
	hurtTime        float64
	currentHurtTime float64
}

func NewEnemy(pos vector) *enemy {
	e := &enemy{
		entity:   NewEntity("enemy-ball-move", pos),
		hurtTime: 0.4 * 1000,
	}
	e.entity.sprites = []*sprite{
		{
			image: "enemy-ball-move",
			pos:   pos,
			animation: &animation{
				numFrames: 4,
				numTime:   0.2 * 1000,
				autoplay:  true,
			},
		},
		{
			image: "enemy-ball-hurt",
			pos:   pos,
			animation: &animation{
				numFrames: 4,
				numTime:   0.1 * 1000,
				autoplay:  true,
			},
		},
	}
	return e
}

func (r *enemy) Update(delta float64) {
	r.entity.Update(delta)
	if r.currentHurtTime > 0 {
		r.currentHurtTime -= delta
		if r.currentHurtTime <= 0 {
			r.entity.SetCurrentSprite(0)
		}
	}
}

func (r *enemy) TakeDamage(amount int) {
	r.entity.health -= amount
	r.currentHurtTime = r.hurtTime
	r.entity.SetCurrentSprite(1)
}
