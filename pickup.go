package raycast

import (
	"fmt"
	"math"
)

type pickupType string

const (
	ammoPickupType   pickupType = "ammo"
	healthPickupType pickupType = "health"
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
	case healthPickupType:
		img = "health"
	}
	return &pickup{
		entity:     NewEntity(img, pos),
		pickupType: t,
		amount:     amount,
	}
}

func (r *pickup) Update(w *World, delta float64) {
	r.entity.Update(delta)
	if r.entity.state != DeadEntityState {
		withinX := math.Abs(w.player.pos.x-r.entity.pos.x) < ((w.player.width + r.entity.width) / 2)
		withinY := math.Abs(w.player.pos.y-r.entity.pos.y) < ((w.player.width + r.entity.width) / 2)
		if withinX && withinY {
			fmt.Printf("entity player collided with entity %v\n ", r.entity.sprite.image)
			switch r.pickupType {
			case ammoPickupType:
				w.player.ammo += r.amount
				if w.player.ammo > maxAmmo {
					w.player.ammo = maxAmmo
				}
			case healthPickupType:
				w.player.health += r.amount
				if w.player.health > maxHealth {
					w.player.health = maxHealth
				}
			}
			r.entity.state = DeadEntityState
		}
	}
}
