package raycast

type effect struct {
	entity *entity
	timer  float64
}

type effectType string

const (
	sceneryDestroyedEffectType = "scenery-destroyed"
	bulletHitEffectType        = "bullet-hit"
	explosionEffectType        = "explosion"
)

func NewEffect(effectType effectType, pos vector) *effect {
	var timing float64
	var numFrames int
	var img string
	switch effectType {
	case bulletHitEffectType:
		timing = 0.08 * 1000
		numFrames = 4
		img = "bullet-hit"
		break
	case sceneryDestroyedEffectType:
		timing = 0.16 * 1000
		numFrames = 4
		img = "grey-hit-effect"
		break
	case explosionEffectType:
		timing = 0.08 * 1000
		numFrames = 8
		img = "explosion"
		break
	}
	e := &effect{
		entity: NewEntity(pos, NewAnimatedSprite(img, &animation{
			numFrames: numFrames,
			numTime:   timing,
			isLoop:    true,
		})),
		timer: float64(numFrames) * timing,
	}
	return e
}

func (r *effect) Update(delta float64, w *World) {
	r.entity.Update(delta, w)
	if r.timer > 0 {
		r.timer -= delta
		if r.timer <= 0 {
			r.entity.state = DeadEntityState
		}
	}
}
