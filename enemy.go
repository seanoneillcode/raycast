package raycast

type enemy struct {
	entity *entity
	// AI state
}

func (r *enemy) Update(delta float64) {
	r.entity.Update(delta)
	// AI behaviour
}
