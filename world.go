package raycast

import "math"

const moveAmount = 0.002
const rotateAmount = 0.0005
const bulletSpeed = 0.01

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
	floorTex   string
	wallTex    string
	ceilingTex string
}

type World struct {
	width   int
	height  int
	tiles   [][]*tile
	bullets []*bullet
	enemies []*enemy
	pickups []*pickup
	effects []*effect
	portals []*portal

	player *player
}

func NewWorld(width, height int) *World {
	w := &World{
		width:  width,
		height: height,
		tiles:  make([][]*tile, width*height),
		player: NewPlayer(vector{
			x: 1.5,
			y: 3,
		}),
		enemies: []*enemy{
			NewEnemy(vector{x: 8, y: 4}),
			NewEnemy(vector{x: 13.5, y: 6}),
			NewEnemy(vector{x: 11, y: 11}),
			NewEnemy(vector{x: 6.5, y: 13}),
			NewEnemy(vector{x: 10.5, y: 7.5}),
			NewEnemy(vector{x: 7.5, y: 9.5}),
		},
		bullets: []*bullet{},
		effects: []*effect{},
		pickups: []*pickup{
			NewPickup(ammoPickupType, 20, vector{x: 10.5, y: 5.5}),
			NewPickup(healthPickupType, 3, vector{x: 8.5, y: 7.5}),
		},
		portals: []*portal{
			NewPortal(vector{x: 2.5, y: 6.5}),
		},
	}
	for x := 0; x < w.width; x++ {
		w.tiles[x] = make([]*tile, width*height)
		for y := 0; y < w.height; y++ {
			w.tiles[x][y] = &tile{}
		}
	}
	initWorld(w)
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

func (w *World) ShootBullet(pos vector, dir vector) {
	w.bullets = append(w.bullets, NewBullet(pos, dir))
}

func (w *World) AddEffect(image string, pos vector) {
	var timing float64
	var numFrames int
	switch image {
	case "bullet-hit":
		timing = 0.08 * 1000
		numFrames = 4
	}
	w.effects = append(w.effects, NewEffect(image, pos, timing, numFrames))
}

func initWorld(w *World) {

	nums := [][]uint8{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 3, 3, 3, 1, 0, 0, 1, 1, 0, 0, 2, 0, 0, 0, 1},
		{1, 3, 3, 3, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 3, 3, 3, 2, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1},
		{1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1},
		{1, 1, 2, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 2, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 3, 3, 3, 3, 3, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 3, 3, 3, 3, 3, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	for iy, y := range nums {
		for ix, x := range y {
			if x == 0 {
				w.tiles[ix][iy].wallTex = "wall"
				w.tiles[ix][iy].floorTex = "floor"
				w.tiles[ix][iy].ceilingTex = "ceiling"
			}
			if x == 1 {
				w.tiles[ix][iy].block = true
				w.tiles[ix][iy].wallTex = "wall"
				w.tiles[ix][iy].floorTex = "floor"
				w.tiles[ix][iy].ceilingTex = "ceiling"
			}
			if x == 2 {
				w.tiles[ix][iy].block = true
				w.tiles[ix][iy].door = true
				w.tiles[ix][iy].wallTex = "door"
				w.tiles[ix][iy].floorTex = "door-floor"
				w.tiles[ix][iy].ceilingTex = "door-floor"
			}
			if x == 3 {
				w.tiles[ix][iy].wallTex = "wall"
				w.tiles[ix][iy].floorTex = "floor"
			}
		}
	}

}
