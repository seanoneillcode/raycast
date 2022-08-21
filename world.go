package raycast

import (
	"math"
	"raycast.com/level"
)

const moveAmount = 0.002
const rotateAmount = 0.0005
const bulletSpeed = 0.008
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
}

type World struct {
	width   int
	height  int
	tiles   [][]*tile
	bullets []*bullet
	enemies []*enemy
	pickups []*pickup
	scenery []*scenery
	effects []*effect
	portals []*portal
	grid    *level.TiledGrid
	player  *player
}

func NewWorld() *World {
	grid := level.NewTileGrid("dungeon.json")
	tiles := loadTiles(grid)
	grid.GetObjectData()
	w := &World{
		width:  grid.Layers[0].Width,
		height: grid.Layers[0].Height,
		tiles:  tiles,
		player: NewPlayer(vector{
			x: 1.5,
			y: 3,
		}),
		enemies: []*enemy{},
		bullets: []*bullet{},
		effects: []*effect{},
		pickups: []*pickup{},
		scenery: []*scenery{},
		grid:    grid,
		portals: []*portal{},
	}
	loadObjectData(grid, w)
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
		b.Update(w, delta)
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

func (w *World) CreatePickup(name string, pos vector) {
	switch name {
	case "soul":
		w.pickups = append(w.pickups, NewPickup(soulPickupType, 1, pos))
		break
	}
}

func loadTiles(grid *level.TiledGrid) [][]*tile {
	tilesRow := make([][]*tile, grid.Layers[0].Width)
	for ix := 0; ix < grid.Layers[0].Width; ix++ {
		tilesColumn := make([]*tile, grid.Layers[0].Height)
		for iy := 0; iy < grid.Layers[0].Height; iy++ {
			td := grid.GetTileData(ix, iy)
			tilesColumn[iy] = &tile{
				block:      td.Block,
				door:       td.Door,
				north:      td.North,
				floorTex:   td.FloorTex,
				wallTex:    td.WallTex,
				doorTex:    td.DoorTex,
				ceilingTex: td.CeilingTex,
			}
		}
		tilesRow[ix] = tilesColumn
	}

	return tilesRow
}

func loadObjectData(grid *level.TiledGrid, w *World) {
	const GridTileSize = 16
	const halfMapTile = 0.5
	objects := grid.GetObjectData()

	for _, obj := range objects {
		pos := vector{
			x: float64(obj.X/GridTileSize) + halfMapTile,
			y: float64(obj.Y/GridTileSize) + halfMapTile,
		}
		switch obj.ObjectType {
		case "level":
			if obj.Name == "start" {
				w.player.pos = pos
			}
			if obj.Name == "end" {
				w.portals = append(w.portals, NewPortal(pos))
			}
			break
		case "scenery":
			if obj.Name == "candlestick" {
				w.scenery = append(w.scenery, NewScenery("candlestick", pos))
			}
			break
		case "enemy":
			if obj.Name == "ball" {
				w.enemies = append(w.enemies, NewEnemy(ballEnemyType, pos))
			}
			if obj.Name == "blue" {
				w.enemies = append(w.enemies, NewEnemy(blueEnemyType, pos))
			}
			break
		case "pickup":
			if obj.Name == "ammo" {
				w.pickups = append(w.pickups, NewPickup(ammoPickupType, 10, pos))
			}
			if obj.Name == "health" {
				w.pickups = append(w.pickups, NewPickup(healthPickupType, 3, pos))
			}
			break
		}

		//for _, p := range obj.Properties {
		//	if p.Name == "team" {
		//		switch p.Value {
		//		case "1":
		//			teamOnePositions = append(teamOnePositions, npc)
		//		case "2":
		//			teamTwoPositions = append(teamTwoPositions, npc)
		//		}
		//	}
		//}

	}
}
