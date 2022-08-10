package raycast

type enemy struct {
	entity            *entity
	hurtTime          float64
	currentHurtTime   float64
	attackTime        float64
	currentAttackTime float64
	state             string
	canSeePlayer      bool
	lastKnowPlayerPos vector
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
			distance: -1,
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
		r.entity.state = HurtEntityState
		if r.currentHurtTime > 0 {
			r.currentHurtTime -= delta
		} else {
			r.entity.SetCurrentSprite(0)
			r.state = "move"
			r.entity.state = NothingEntityState
		}
		break
	case "move":
		if r.entity.CurrentSprite().distance != -1 && r.entity.CurrentSprite().distance < 1.0 {
			r.entity.SetCurrentSprite(2)
			r.state = "attack"
			r.currentAttackTime = r.attackTime
		}

		break
	case "attack":
		if r.currentAttackTime > 0 {
			r.currentAttackTime -= delta
		} else {
			if r.entity.CurrentSprite().distance < 1.0 {
				w.player.TakeDamage(1)
			}
			r.entity.SetCurrentSprite(0)
			r.state = "move"
		}
		r.entity.state = NothingEntityState
		break
	}
	if r.canSeePlayer {
		if !canSeePos(w, r.entity.pos, w.player.pos) {
			r.canSeePlayer = false
			r.entity.dir = normalizeVector(vector{
				x: r.lastKnowPlayerPos.x - r.entity.pos.x,
				y: r.lastKnowPlayerPos.y - r.entity.pos.y,
			})
		} else {
			r.lastKnowPlayerPos.x = w.player.pos.x
			r.lastKnowPlayerPos.y = w.player.pos.y
		}
		r.entity.dir = normalizeVector(vector{
			x: w.player.pos.x - r.entity.pos.x,
			y: w.player.pos.y - r.entity.pos.y,
		})
	} else {
		if within(r.entity.pos, r.lastKnowPlayerPos, 0.25) {
			r.entity.dir.x = 0
			r.entity.dir.y = 0
		}
		if canSeePos(w, r.entity.pos, w.player.pos) {
			r.canSeePlayer = true
		}
	}
}

func (r *enemy) TakeDamage(amount int) {
	r.entity.health -= amount
	r.currentHurtTime = r.hurtTime
	r.entity.SetCurrentSprite(1)
	r.state = "hurt"
}
