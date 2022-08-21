package raycast

import (
	"math"
)

type ray struct {
	side     int
	distance float64
	wallX    float64
	dir      vector
	texture  string
}

func calculateRay(w *World, cameraX float64) ray {
	rayStart := vector{
		x: w.player.pos.x,
		y: w.player.pos.y,
	}

	rayDir := vector{
		x: w.player.dir.x + w.player.plane.x*cameraX,
		y: w.player.dir.y + w.player.plane.y*cameraX,
	}

	rayUnitStepSize := vector{
		x: math.Abs(1 / rayDir.x),
		y: math.Abs(1 / rayDir.y),
	}

	if rayUnitStepSize.x == 0 {
		rayUnitStepSize.x = 111111111
	}
	if rayUnitStepSize.y == 0 {
		rayUnitStepSize.y = 111111111
	}

	rayMapPos := mapPos{
		x: int(rayStart.x),
		y: int(rayStart.y),
	}

	rayLength := vector{}
	step := mapPos{}

	if rayDir.x < 0 {
		step.x = -1
		rayLength.x = (rayStart.x - float64(rayMapPos.x)) * rayUnitStepSize.x
	} else {
		step.x = 1
		rayLength.x = (float64(rayMapPos.x+1) - rayStart.x) * rayUnitStepSize.x
	}
	if rayDir.y < 0 {
		step.y = -1
		rayLength.y = (rayStart.y - float64(rayMapPos.y)) * rayUnitStepSize.y
	} else {
		step.y = 1
		rayLength.y = (float64(rayMapPos.y+1) - rayStart.y) * rayUnitStepSize.y
	}

	tileFound := false
	maxDistance := 256.0
	distance := 0.0

	var texture string
	var side = 0
	var t *tile

	for !tileFound && distance < maxDistance {

		t = w.getTile(rayMapPos.x, rayMapPos.y)
		if t != nil {
			texture = t.wallTex
			if t.block {
				tileFound = true
				if t.door {
					texture = t.doorTex
					if t.north {
						if rayLength.y > rayLength.x && rayUnitStepSize.y < 0.5 {
							side = 0
							texture = t.doorTex
						} else {
							rayLength.y = rayLength.y - (rayUnitStepSize.y / 2)

						}
					} else {
						if rayLength.x > rayLength.y && rayUnitStepSize.x < 0.5 {
							side = 1
							texture = t.doorTex
						} else {
							rayLength.x = rayLength.x - (rayUnitStepSize.x / 2)
						}
					}
				}
			}
			if t.door {
				if t.north {
					if rayLength.y > rayLength.x {
						tileFound = true
						side = 0
						texture = t.wallTex
					}
				} else {
					if rayLength.x > rayLength.y {
						tileFound = true
						side = 1
						texture = t.wallTex
					}
				}
			}
			t.seen = true
		}
		if !tileFound {
			if rayLength.x < rayLength.y {
				rayMapPos.x += step.x
				distance = rayLength.x
				rayLength.x += rayUnitStepSize.x
				side = 0
			} else {
				rayMapPos.y += step.y
				distance = rayLength.y
				rayLength.y += rayUnitStepSize.y
				side = 1
			}
		}
	}

	perpWallDist := 256.0
	if tileFound {
		if t != nil && t.door {
			if side == 0 {
				perpWallDist = rayLength.x
			} else {
				perpWallDist = rayLength.y
			}
		} else {
			if side == 0 {
				perpWallDist = rayLength.x - rayUnitStepSize.x
			} else {
				perpWallDist = rayLength.y - rayUnitStepSize.y
			}
		}
	}

	var wallX float64
	if side == 0 {
		wallX = rayStart.y + (perpWallDist * rayDir.y)
	} else {
		wallX = rayStart.x + (perpWallDist * rayDir.x)
	}
	wallX -= math.Floor(wallX)

	return ray{
		distance: perpWallDist,
		side:     side,
		wallX:    wallX,
		dir:      rayDir,
		texture:  texture,
	}
}

func canSeePos(w *World, startPos vector, targetPos vector) bool {
	rayStart := startPos
	targetMapPos := mapPos{
		x: int(targetPos.x),
		y: int(targetPos.y),
	}

	rayDir := normalizeVector(vector{
		x: targetPos.x - startPos.x,
		y: targetPos.y - startPos.y,
	})

	rayUnitStepSize := vector{
		x: math.Abs(1 / rayDir.x),
		y: math.Abs(1 / rayDir.y),
	}

	if rayUnitStepSize.x == 0 {
		rayUnitStepSize.x = 111111111
	}
	if rayUnitStepSize.y == 0 {
		rayUnitStepSize.y = 111111111
	}

	rayMapPos := mapPos{
		x: int(rayStart.x),
		y: int(rayStart.y),
	}

	rayLength := vector{}
	step := mapPos{}

	if rayDir.x < 0 {
		step.x = -1
		rayLength.x = (rayStart.x - float64(rayMapPos.x)) * rayUnitStepSize.x
	} else {
		step.x = 1
		rayLength.x = (float64(rayMapPos.x+1) - rayStart.x) * rayUnitStepSize.x
	}
	if rayDir.y < 0 {
		step.y = -1
		rayLength.y = (rayStart.y - float64(rayMapPos.y)) * rayUnitStepSize.y
	} else {
		step.y = 1
		rayLength.y = (float64(rayMapPos.y+1) - rayStart.y) * rayUnitStepSize.y
	}

	maxDistance := 256.0
	distance := 0.0

	for distance < maxDistance {

		if rayLength.x < rayLength.y {
			rayMapPos.x += step.x
			distance += rayLength.x
			rayLength.x += rayUnitStepSize.x
		} else {
			rayMapPos.y += step.y
			distance += rayLength.y
			rayLength.y += rayUnitStepSize.y
		}

		t := w.getTile(rayMapPos.x, rayMapPos.y)
		if t != nil {
			if t.block {
				return false
			}
		}
		if rayMapPos.x == targetMapPos.x && rayMapPos.y == targetMapPos.y {
			return true
		}
	}

	return false
}
