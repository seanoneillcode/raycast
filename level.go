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
			}
		}
		tilesRow[ix] = tilesColumn
	}

	return tilesRow
}

type objectData struct {
	startPos vector
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
			}
			if obj.Name == "end" {
				objData.portals = append(objData.portals, NewPortal(pos))
			}
			break
		case "scenery":
			if obj.Name == "candlestick" {
				objData.scenery = append(objData.scenery, NewScenery("candlestick", pos))
			}
			break
		case "enemy":
			if obj.Name == "ball" {
				objData.enemies = append(objData.enemies, NewEnemy(ballEnemyType, pos))
			}
			if obj.Name == "blue" {
				objData.enemies = append(objData.enemies, NewEnemy(blueEnemyType, pos))
			}
			break
		case "pickup":
			if obj.Name == "ammo" {
				objData.pickups = append(objData.pickups, NewPickup(ammoPickupType, 10, pos))
			}
			if obj.Name == "health" {
				objData.pickups = append(objData.pickups, NewPickup(healthPickupType, 3, pos))
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
