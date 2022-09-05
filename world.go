package raycast

import (
	"math"
)

const moveAmount = 0.002
const rotateAmount = 0.0005
const bulletSpeed = 0.01
const bulletWidth = 0.01

type vector struct {
	x float64
	y float64
}

func addVector(v1, v2 vector) vector {
	return vector{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
	}
}

func scaleVector(v1 vector, amount float64) vector {
	return vector{
		x: v1.x * amount,
		y: v1.y * amount,
	}
}

func normalizeVector(v1 vector) vector {
	magnitude := math.Sqrt((v1.x * v1.x) + (v1.y * v1.y))
	return vector{
		x: v1.x / magnitude,
		y: v1.y / magnitude,
	}
}

type mapPos struct {
	x int
	y int
}

type tile struct {
	block      bool
	door       bool
	north      bool
	floorTex   string
	wallTex    string
	doorTex    string
	ceilingTex string
	seen       bool
	locked     bool
}

type World struct {
	width       int
	height      int
	tiles       [][]*tile
	bullets     []*bullet
	enemies     []*enemy
	pickups     []*pickup
	scenery     []*scenery
	effects     []*effect
	portals     []*portal
	player      *player
	soundPlayer *SoundPlayer
}

func NewWorld() *World {
	l := LoadLevel("library.json")

	w := &World{
		soundPlayer: NewSoundPlayer(),
		tiles:       l.tiles,
		width:       l.width,
		height:      l.height,
		player:      NewPlayer(l.objectData.startPos, l.objectData.startDir),
		enemies:     l.objectData.enemies,
		pickups:     l.objectData.pickups,
		scenery:     l.objectData.scenery,
		portals:     l.objectData.portals,
		// temp state
		bullets: []*bullet{},
		effects: []*effect{},
	}
	w.soundPlayer.LoadSound("pickup-health")
	w.soundPlayer.LoadSound("pickup-ammo")
	w.soundPlayer.LoadSound("pickup-soul")
	w.soundPlayer.LoadSound("door")
	w.soundPlayer.LoadSound("crack")
	w.soundPlayer.LoadSound("thud")
	w.soundPlayer.LoadSound("chunk")
	w.soundPlayer.LoadSound("player-hurt")
	w.soundPlayer.LoadSound("bullet-hit")
	w.soundPlayer.LoadSound("enemy-die")
	w.soundPlayer.LoadSound("enemy-hurt")
	w.soundPlayer.LoadSound("enemy-shoot")
	return w
}

func (w *World) Update(delta float64) error {
	hasDead := false
	for _, e := range w.enemies {
		e.Update(w, delta)
		if e.entity.state == DeadEntityState {
			hasDead = true
		}
	}
	if hasDead {
		cleanDeadEnemy(w)
	}

	hasDead = false
	for _, b := range w.bullets {
		b.Update(w, delta)
		if b.entity.state == DeadEntityState {
			hasDead = true
		}
	}
	if hasDead {
		cleanDeadBullet(w)
	}

	hasDead = false
	for _, b := range w.pickups {
		b.Update(w, delta)
		if b.entity.state == DeadEntityState {
			hasDead = true
		}
	}
	if hasDead {
		cleanDeadPickup(w)
	}

	hasDead = false
	for _, b := range w.scenery {
		b.Update(delta)
		if b.entity.state == DeadEntityState {
			hasDead = true
		}
	}
	if hasDead {
		cleanDeadScenery(w)
	}

	hasDead = false
	for _, b := range w.effects {
		b.Update(delta)
		if b.entity.state == DeadEntityState {
			hasDead = true
		}
	}
	if hasDead {
		cleanDeadEffect(w)
	}

	for _, b := range w.portals {
		b.Update(w, delta)
	}

	err := w.player.Update(w, delta)
	if err != nil {
		return err
	}

	return nil
}

func cleanDeadScenery(w *World) {
	temp := w.scenery[:0]
	for _, b := range w.scenery {
		if b.entity.state != DeadEntityState {
			temp = append(temp, b)
		}
	}
	w.scenery = temp
}

func cleanDeadPickup(w *World) {
	temp := w.pickups[:0]
	for _, b := range w.pickups {
		if b.entity.state != DeadEntityState {
			temp = append(temp, b)
		}
	}
	w.pickups = temp
}

func cleanDeadEffect(w *World) {
	temp := w.effects[:0]
	for _, b := range w.effects {
		if b.entity.state != DeadEntityState {
			temp = append(temp, b)
		}
	}
	w.effects = temp
}

func cleanDeadBullet(w *World) {
	temp := w.bullets[:0]
	for _, b := range w.bullets {
		if b.entity.state != DeadEntityState {
			temp = append(temp, b)
		}
	}
	w.bullets = temp
}

func cleanDeadEnemy(w *World) {
	temp := w.enemies[:0]
	for _, b := range w.enemies {
		if b.entity.state != DeadEntityState {
			temp = append(temp, b)
		}
	}
	w.enemies = temp
}

func (w *World) getTileAtPoint(pos vector) *tile {
	return w.getTile(int(pos.x), int(pos.y))
}

func (w *World) getTile(x, y int) *tile {
	if x < 0 || x > w.width-1 {
		return nil
	}
	if y < 0 || y > w.height-1 {
		return nil
	}
	return w.tiles[x][y]
}

func (w *World) ShootBullet(pos vector, dir vector, speed float64) {
	w.bullets = append(w.bullets, NewBullet(pos, dir, speed))
}

func (w *World) AddEffect(effectType effectType, pos vector) {
	w.effects = append(w.effects, NewEffect(effectType, pos))
}

func (w *World) CreateEntity(name string, pos vector) {
	switch name {
	case "soul":
		w.pickups = append(w.pickups, NewPickup(soulPickupType, 1, pos))
	case "health":
		w.pickups = append(w.pickups, NewPickup(healthPickupType, 1, pos))
	case "ammo":
		w.pickups = append(w.pickups, NewPickup(ammoPickupType, 5, pos))
	case "end":
		w.portals = append(w.portals, NewPortal(pos))
	}
}
