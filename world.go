package raycast

import "fmt"

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

type mapPos struct {
	x int
	y int
}

type tile struct {
	block bool
}

type World struct {
	width   int
	height  int
	tiles   [][]*tile
	bullets []*bullet
	enemies []*enemy

	player *player
}

func NewWorld(width, height int) *World {
	w := &World{
		width:  width,
		height: height,
		tiles:  make([][]*tile, width*height),
		player: NewPlayer(vector{
			x: 5,
			y: 6,
		}),
		enemies: []*enemy{
			{
				entity: NewEntity("eye", vector{x: 5, y: 6}),
			},
		},
		bullets: []*bullet{
			{
				entity: NewEntity("bullet", vector{x: 10, y: 6}),
			},
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
		e.Update(delta)
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
	err := w.player.Update(w, delta)
	if err != nil {
		return err
	}

	return nil
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
		} else {
			fmt.Printf("removing dead entity: %v", b.entity.sprite)
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
	b := &bullet{
		entity: NewEntity("bullet", pos),
	}
	b.entity.dir = dir
	b.entity.speed = bulletSpeed
	w.bullets = append(w.bullets, b)
	fmt.Printf("number of bullets in world: %v \n", len(w.bullets))
	for i, b := range w.bullets {
		fmt.Printf("bullet %v pos %v\n", i, b.entity.pos)
	}
}

func initWorld(w *World) {

	nums := [][]uint8{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1},
		{1, 1, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1},
		{1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1},
		{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	for ix, x := range nums {
		for iy, y := range x {
			if y == 1 {
				w.tiles[ix][iy].block = true
			}
		}
	}

}
