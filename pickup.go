package raycast

import (
	"math"
)

type pickupType string

const (
	ammoPickupType   pickupType = "ammo"
	healthPickupType pickupType = "health"
	soulPickupType   pickupType = "soul"
)

type pickup struct {
	entity     *entity
	pickupType pickupType
	amount     int
}

func NewPickup(t pickupType, amount int, pos vector) *pickup {
	var img string
	switch t {
	case ammoPickupType:
		img = "ammo"
		break
	case healthPickupType:
		img = "health"
		break
	case soulPickupType:
		img = "soul"
		break
	}
	p := &pickup{
		entity:     NewEntity(img, pos),
		pickupType: t,
		amount:     amount,
	}
	if t == soulPickupType {
		p.entity.CurrentSprite().animation = &animation{
			numFrames: 4,
			numTime:   0.2 * 1000,
			autoplay:  true,
		}
	}
	return p
}

func (r *pickup) Update(w *World, delta float64) {
	r.entity.Update(delta)
	if r.entity.state != DeadEntityState {
		withinX := math.Abs(w.player.pos.x-r.entity.pos.x) < ((w.player.width + r.entity.width) / 2)
		withinY := math.Abs(w.player.pos.y-r.entity.pos.y) < ((w.player.width + r.entity.width) / 2)
		if withinX && withinY {
			switch r.pickupType {
			case ammoPickupType:
				w.player.ammo += r.amount
				if w.player.ammo > maxAmmo {
					w.player.ammo = maxAmmo
				}
				w.player.screenFlashTimer = screenFlashTime
				w.player.screenFlashColor = ammoPickupScreenFlashColor
				break
			case healthPickupType:
				w.player.health += r.amount
				if w.player.health > maxHealth {
					w.player.health = maxHealth
				}
				w.player.screenFlashTimer = screenFlashTime
				w.player.screenFlashColor = healthPickupScreenFlashColor
				break
			case soulPickupType:
				w.player.souls += r.amount
				w.player.screenFlashTimer = screenFlashTime
				w.player.screenFlashColor = soulPickupScreenFlashColor
				break
			}
			r.entity.state = DeadEntityState
		}
	}
}
