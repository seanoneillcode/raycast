package raycast

const moveAmount = 0.002
const rotateAmount = 0.0005

type vector struct {
	x float64
	y float64
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

	for _, e := range w.enemies {
		e.Update(delta)
	}
	for _, b := range w.bullets {
		b.Update(delta)
	}
	err := w.player.Update(w, delta)
	if err != nil {
		return err
	}

	return nil
}

func (w *World) getTileAtPoint(x, y float64) *tile {
	return w.getTile(int(x), int(y))
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
