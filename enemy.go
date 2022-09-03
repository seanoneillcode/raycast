package raycast

type enemy struct {
	entity            *entity
	hurtTime          float64
	dyingTime         float64
	currentHurtTime   float64
	attackTime        float64
	currentAttackTime float64
	state             string
	canSeePlayer      bool
	lastKnowPlayerPos vector
	enemyType         EnemyType
	attackRange       float64
}

type EnemyType string

const (
	blueEnemyType = "blue-enemy"
	ballEnemyType = "ball-enemy"
)

func NewEnemy(enemyType EnemyType, pos vector) *enemy {
	e := &enemy{
		hurtTime:    0.6 * 1000,
		dyingTime:   0.4 * 1000,
		attackTime:  0.6 * 1000,
		enemyType:   enemyType,
		attackRange: 1.0,
		state:       "move",
	}
	var ent *entity
	switch enemyType {
	case ballEnemyType:
		ent = NewEntity(
			pos,
			NewAnimatedSprite("enemy-ball-move", &animation{
				numFrames: 4,
				numTime:   0.2 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-ball-hurt", &animation{
				numFrames: 4,
				numTime:   0.15 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-ball-attack", &animation{
				numFrames: 4,
				numTime:   0.15 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-ball-die", &animation{
				numFrames: 4,
				numTime:   0.1 * 1000,
				autoplay:  true,
			}),
		)
		ent.speed = 0.003
		ent.dropItem = "soul"
		break
	case blueEnemyType:
		e.attackRange = 8 * 8 // distance needs to be squared

		ent = NewEntity(
			pos,
			NewAnimatedSprite("enemy-blue-move", &animation{
				numFrames: 4,
				numTime:   0.2 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-blue-hurt", &animation{
				numFrames: 4,
				numTime:   0.15 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-blue-attack", &animation{
				numFrames: 4,
				numTime:   0.15 * 1000,
				autoplay:  true,
			}),
			NewAnimatedSprite("enemy-blue-die", &animation{
				numFrames: 4,
				numTime:   0.1 * 1000,
				autoplay:  true,
			}),
		)
		ent.speed = 0.001
		ent.dropItem = "soul"
		break
	}
	e.entity = ent

	return e
}

func (r *enemy) Update(w *World, delta float64) {
	r.entity.Update(delta)
	if r.entity.health < 0 && r.state != "dying" {
		w.soundPlayer.PlaySound("enemy-die")
		r.state = "dying"
		r.currentHurtTime = r.dyingTime
		r.entity.SetCurrentSprite(3)
		if r.entity.dropItem != "" {
			w.CreateEntity(r.entity.dropItem, r.entity.pos)
		}
	}

	switch r.state {
	case "hurt":
		r.entity.state = StunnedEntityState
		if r.currentHurtTime > 0 {
			r.currentHurtTime -= delta
		} else {

			r.entity.SetCurrentSprite(0)
			r.state = "move"
			r.entity.state = NothingEntityState
		}
		break
	case "move":
		if r.entity.CurrentSprite().distance != -1 && r.entity.CurrentSprite().distance < r.attackRange {
			r.entity.SetCurrentSprite(2)
			r.state = "attack"
			r.currentAttackTime = r.attackTime
		}
		break
	case "attack":
		switch r.enemyType {
		case ballEnemyType:
			if r.currentAttackTime > 0 {
				r.currentAttackTime -= delta
			} else {
				if r.entity.CurrentSprite().distance < r.attackRange {
					w.player.TakeDamage(1)
					w.soundPlayer.PlaySound("enemy-shoot")
				}
				r.entity.SetCurrentSprite(0)
				r.state = "move"
			}
			break
		case blueEnemyType:
			if r.currentAttackTime > 0 {
				r.currentAttackTime -= delta
			} else {
				if canSeePos(w, r.entity.pos, w.player.pos) {
					bulletDir := normalizeVector(vector{
						x: w.player.pos.x - r.entity.pos.x,
						y: w.player.pos.y - r.entity.pos.y,
					})
					w.ShootBullet(addVector(r.entity.pos, bulletDir), bulletDir, bulletSpeed/2)
					w.soundPlayer.PlaySound("enemy-shoot")
				}
				r.entity.SetCurrentSprite(0)
				r.state = "move"
			}
			break
		}
		r.entity.state = NothingEntityState
		break
	case "dying":
		r.entity.state = StunnedEntityState
		if r.currentHurtTime > 0 {
			r.currentHurtTime -= delta
		} else {
			r.entity.state = DeadEntityState
		}
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

func (r *enemy) TakeDamage(w *World, amount int) {
	if r.state == "dying" {
		return
	}
	r.entity.health -= amount
	r.currentHurtTime = r.hurtTime
	r.entity.SetCurrentSprite(1)
	r.state = "hurt"
	w.soundPlayer.PlaySound("chunk")
}
