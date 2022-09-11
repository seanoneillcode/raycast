package raycast

type bullet struct {
	entity *entity
}

func NewBullet(pos vector, dir vector, speed float64) *bullet {
	b := &bullet{
		entity: NewEntity(pos, NewSprite("bullet")),
	}
	b.entity.dir = dir
	b.entity.speed = speed
	b.entity.width = bulletWidth
	return b
}

func (r *bullet) Update(w *World, delta float64) {
	r.entity.Update(delta, w)
	t := w.getTileAtPoint(r.entity.pos)
	if t.block {
		r.entity.state = DeadEntityState
		r.entity.undoLastMove(delta)
		w.AddEffect(bulletHitEffectType, r.entity.pos)
		w.soundPlayer.PlaySound("thud")
	}
	for _, e := range w.enemies {
		if collides(r.entity, e.entity) {
			r.entity.state = DeadEntityState
			r.entity.undoLastMove(delta)
			w.AddEffect(bulletHitEffectType, r.entity.pos)
			e.TakeDamage(w, 1)
			w.soundPlayer.PlaySound("thud")
		}
	}
	for _, e := range w.scenery {
		if e.canCollide && collides(r.entity, e.entity) {
			r.entity.state = DeadEntityState
			r.entity.undoLastMove(delta)
			w.AddEffect(bulletHitEffectType, r.entity.pos)
			e.TakeDamage(w, 1)
			w.soundPlayer.PlaySound("thud")
		}
	}
	if collidesWithPlayer(w.player, r.entity) {
		r.entity.state = DeadEntityState
		r.entity.undoLastMove(delta)
		w.player.TakeDamage(1)
	}
}
