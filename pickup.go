package raycast

import (
	"math"
)

type pickupType string

const (
	ammoPickupType    pickupType = "ammo"
	healthPickupType  pickupType = "health"
	soulPickupType    pickupType = "soul"
	crystalPickupType pickupType = "crystal"
)

type pickup struct {
	entity     *entity
	pickupType pickupType
	amount     int
}

func NewPickup(t pickupType, amount int, pos vector) *pickup {
	var s *sprite
	switch t {
	case ammoPickupType:
		s = NewSprite("ammo")
		break
	case healthPickupType:
		s = NewSprite("health")
		break
	case soulPickupType:
		s = NewAnimatedSprite("soul", &animation{
			numFrames: 4,
			numTime:   0.2 * 1000,
			autoplay:  true,
		})
		break
	case crystalPickupType:
		s = NewAnimatedSprite("book", &animation{
			numFrames: 4,
			numTime:   0.4 * 1000,
			autoplay:  true,
		})
		break
	}
	p := &pickup{
		entity:     NewEntity(pos, s),
		pickupType: t,
		amount:     amount,
	}
	return p
}

func (r *pickup) Update(w *World, delta float64) {
	r.entity.Update(delta)
	if r.entity.state != DeadEntityState {
		withinX := math.Abs(w.player.pos.x-r.entity.pos.x) < ((w.player.width + r.entity.width) / 2)
		withinY := math.Abs(w.player.pos.y-r.entity.pos.y) < ((w.player.width + r.entity.width) / 2)
		if withinX && withinY {
			handleGettingPickedUp(w, r)
			r.entity.state = DeadEntityState
		}
	}
}

func handleGettingPickedUp(w *World, p *pickup) {
	switch p.pickupType {
	case ammoPickupType:
		w.player.ammo += p.amount
		if w.player.ammo > maxAmmo {
			w.player.ammo = maxAmmo
		}
		w.player.screenFlashTimer = screenFlashTime
		w.player.screenFlashColor = ammoPickupScreenFlashColor
		w.soundPlayer.PlaySound("pickup-ammo")
		break
	case healthPickupType:
		w.player.health += p.amount
		if w.player.health > maxHealth {
			w.player.health = maxHealth
		}
		w.player.screenFlashTimer = screenFlashTime
		w.player.screenFlashColor = healthPickupScreenFlashColor
		w.soundPlayer.PlaySound("pickup-health")
		break
	case soulPickupType:
		w.player.souls += p.amount
		w.player.screenFlashTimer = screenFlashTime
		w.player.screenFlashColor = soulPickupScreenFlashColor
		w.soundPlayer.PlaySound("pickup-soul")
		break
	case crystalPickupType:
		w.player.souls += p.amount
		w.player.screenFlashTimer = screenFlashTime
		w.player.screenFlashColor = soulPickupScreenFlashColor
		w.soundPlayer.PlaySound("pickup-soul")
		break
	}
}
