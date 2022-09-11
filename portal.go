package raycast

import (
	"fmt"
	"math"
)

type portal struct {
	entity *entity
}

func NewPortal(pos vector) *portal {
	timing := 0.2 * 1000
	p := &portal{
		entity: NewEntity(pos, NewAnimatedSprite("portal", &animation{
			numFrames: 4,
			numTime:   timing,
			autoplay:  true,
		})),
	}
	return p
}

func (r *portal) Update(w *World, delta float64) {
	r.entity.Update(delta, w)
	withinX := math.Abs(w.player.pos.x-r.entity.pos.x) < ((w.player.width + r.entity.width) / 2)
	withinY := math.Abs(w.player.pos.y-r.entity.pos.y) < ((w.player.width + r.entity.width) / 2)
	if withinX && withinY {
		fmt.Printf("entity player collided with portal\n ")
		fmt.Printf("level won!!")
		panic("player won")
	}
}
