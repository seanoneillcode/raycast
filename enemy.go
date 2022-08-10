package raycast

type enemy struct {
	entity            *entity
	hurtTime          float64
	currentHurtTime   float64
	attackTime        float64
	currentAttackTime float64
	state             string
}

func NewEnemy(pos vector) *enemy {
	e := &enemy{
		entity:     NewEntity("enemy-ball-move", pos),
		hurtTime:   0.6 * 1000,
		attackTime: 0.6 * 1000,
		state:      "move",
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
				numTime:   0.15 * 1000,
				autoplay:  true,
			},
		},
		{
			image: "enemy-ball-attack",
			pos:   pos,
			animation: &animation{
				numFrames: 4,
				numTime:   0.15 * 1000,
				autoplay:  true,
			},
		},
	}
	return e
}

func (r *enemy) Update(w *World, delta float64) {
	r.entity.Update(delta)

	switch r.state {
	case "hurt":
		if r.currentHurtTime > 0 {
			r.currentHurtTime -= delta
		} else {
			r.entity.SetCurrentSprite(0)
			r.state = "move"
		}
	case "move":
		if r.entity.CurrentSprite().distance < 4.0 {
			r.entity.SetCurrentSprite(2)
			r.state = "attack"
			r.currentAttackTime = r.attackTime
		}
	case "attack":
		if r.currentAttackTime > 0 {
			r.currentAttackTime -= delta
		} else {
			if r.entity.CurrentSprite().distance < 2.0 {
				w.player.TakeDamage(1)
			}
			r.entity.SetCurrentSprite(0)
			r.state = "move"
		}
	}
}

func (r *enemy) TakeDamage(amount int) {
	r.entity.health -= amount
	r.currentHurtTime = r.hurtTime
	r.entity.SetCurrentSprite(1)
	r.state = "hurt"
}
