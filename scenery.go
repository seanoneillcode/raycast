package raycast

type scenery struct {
	entity *entity
}

func NewScenery(img string, pos vector) *scenery {
	p := &scenery{
		entity: NewEntity(img, pos),
	}
	p.entity.health = 0
	if img == "candlestick" {
		p.entity.CurrentSprite().animation = &animation{
			numFrames: 4,
			numTime:   0.2 * 1000,
			autoplay:  true,
		}
	}
	return p
}

func (r *scenery) Update(w *World, delta float64) {
	r.entity.Update(delta)
}

func (r *scenery) TakeDamage(w *World, amount int) {
	r.entity.health -= amount
	if r.entity.health < 0 {
		r.entity.state = DeadEntityState
		w.AddEffect("grey-hit-effect", r.entity.pos)
		w.soundPlayer.PlaySound("enemy-hurt")
	}
}
