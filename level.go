package raycast

import "raycast.com/tiledgrid"

type level struct {
	objectData *objectData
	tiles      [][]*tile
	width      int
	height     int
}

func LoadLevel(fileName string) *level {
	grid := tiledgrid.NewTileGrid(fileName)
	return &level{
		tiles:      loadTiles(grid),
		objectData: loadObjectData(grid),
		width:      grid.Layers[0].Width,
		height:     grid.Layers[0].Height,
	}
}

func loadTiles(grid *tiledgrid.TiledGrid) [][]*tile {
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
				locked:     td.Locked,
			}
		}
		tilesRow[ix] = tilesColumn
	}

	return tilesRow
}

type objectData struct {
	startPos vector
	startDir string
	enemies  []*enemy
	pickups  []*pickup
	scenery  []*scenery
	portals  []*portal
}

func loadObjectData(grid *tiledgrid.TiledGrid) *objectData {
	objData := &objectData{
		enemies: []*enemy{},
		pickups: []*pickup{},
		scenery: []*scenery{},
		portals: []*portal{},
	}

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
				objData.startPos = pos
				objData.startDir = getStringProperty("dir", obj)
			}
			if obj.Name == "end" {
				objData.portals = append(objData.portals, NewPortal(pos))
			}
			break
		case "scenery":
			if obj.Name == "candlestick" {
				s := NewAnimatedSprite("candlestick", &animation{
					numFrames: 4,
					numTime:   0.2 * 1000,
					autoplay:  true,
				})
				objData.scenery = append(objData.scenery, NewScenery(s, pos, sceneryDestroyedEffectType, "enemy-hurt", "", true))
			}
			if obj.Name == "barrel" {
				s := NewSprite("barrel")
				s.height = 0.5
				objData.scenery = append(objData.scenery, NewScenery(s, pos, explosionEffectType, "enemy-die", "", true))
			}
			if obj.Name == "tree" {
				s := NewSprite("tree")
				objData.scenery = append(objData.scenery, NewScenery(s, pos, explosionEffectType, "thud", "", false))
			}
			if obj.Name == "bush" {
				s := NewSprite("bush")
				objData.scenery = append(objData.scenery, NewScenery(s, pos, explosionEffectType, "thud", "", false))
			}
			if obj.Name == "web" {
				s := NewSprite("web")
				objData.scenery = append(objData.scenery, NewScenery(s, pos, sceneryDestroyedEffectType, "enemy-hurt", "", false))
			}
			break
		case "enemy":
			if obj.Name == "ball" {
				objData.enemies = append(objData.enemies, NewEnemy(ballEnemyType, pos))
			}
			if obj.Name == "blue" {
				objData.enemies = append(objData.enemies, NewEnemy(blueEnemyType, pos))
			}
			if obj.Name == "blob" {
				objData.enemies = append(objData.enemies, NewEnemy(blobEnemyType, pos))
			}
			if obj.Name == "alien" {
				objData.enemies = append(objData.enemies, NewEnemy(alienEnemyType, pos))
			}
			break
		case "pickup":
			if obj.Name == "ammo" {
				objData.pickups = append(objData.pickups, NewPickup(ammoPickupType, 10, pos))
			}
			if obj.Name == "health" {
				objData.pickups = append(objData.pickups, NewPickup(healthPickupType, 3, pos))
			}
			if obj.Name == "book" {
				objData.pickups = append(objData.pickups, NewPickup(bookPickupType, 1, pos))
			}
			if obj.Name == "key" {
				objData.pickups = append(objData.pickups, NewPickup(keyPickupType, 1, pos))
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
	return objData
}

func getStringProperty(name string, obj *tiledgrid.ObjectData) string {
	for _, p := range obj.Properties {
		if p.Name == name {
			return (p.Value).(string)
		}
	}
	return ""
}
