package raycast

type effect struct {
	entity *entity
	timer  float64
}

func NewEffect(image string, pos vector, timing float64, numFrames int) *effect {
	e := &effect{
		entity: NewEntity(image, pos),
		timer:  float64(numFrames) * timing,
	}
	e.entity.sprites[0].animation = &animation{
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
