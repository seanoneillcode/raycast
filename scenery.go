package raycast

type scenery struct {
	entity *entity
	effect effectType
	sound  string
}

func NewScenery(sprite *sprite, pos vector, effect effectType, sound string, dropItem string, isPhysicsEntity bool) *scenery {
	p := &scenery{
		entity: NewEntity(pos, sprite),
		effect: effect,
		sound:  sound,
	}
	p.entity.isPhysicsEntity = isPhysicsEntity
	if dropItem != "" {
		p.entity.dropItem = dropItem
	}
	p.entity.health = 0
	return p
}

func (r *scenery) Update(delta float64, w *World) {
	r.entity.Update(delta, w)
}

func (r *scenery) TakeDamage(w *World, amount int) {
	r.entity.health -= amount
	if r.entity.health < 0 {
		r.entity.state = DeadEntityState
		w.AddEffect(r.effect, r.entity.pos)
		w.AddParticles(smokeParticleType, r.entity.pos)
		w.soundPlayer.PlaySound(r.sound)
		if r.entity.dropItem != "" {
			w.CreateEntity(r.entity.dropItem, r.entity.pos)
		}
	}
}
