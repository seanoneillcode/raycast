package raycast

import "math"

type ray struct {
	side     int
	distance float64
	wallx    float64
	dir      point
}

func calculateRay(w *World, cameraX float64) ray {
	rayStart := point{
		x: w.playerPos.x,
		y: w.playerPos.y,
	}

	rayDir := point{
		x: w.playerDir.x + w.plane.x*cameraX,
		y: w.playerDir.y + w.plane.y*cameraX,
	}

	rayUnitStepSize := point{
		x: math.Abs(1 / rayDir.x),
		y: math.Abs(1 / rayDir.y),
	}

	if rayUnitStepSize.x == 0 {
		rayUnitStepSize.x = 111111111
	}
	if rayUnitStepSize.y == 0 {
		rayUnitStepSize.y = 111111111
	}

	mapPos := pos{
		x: int(rayStart.x),
		y: int(rayStart.y),
	}

	rayLength := point{}
	Step := pos{}

	if rayDir.x < 0 {
		Step.x = -1
		rayLength.x = (rayStart.x - float64(mapPos.x)) * rayUnitStepSize.x
	} else {
		Step.x = 1
		rayLength.x = (float64(mapPos.x+1) - rayStart.x) * rayUnitStepSize.x
	}
	if rayDir.y < 0 {
		Step.y = -1
		rayLength.y = (rayStart.y - float64(mapPos.y)) * rayUnitStepSize.y
	} else {
		Step.y = 1
		rayLength.y = (float64(mapPos.y+1) - rayStart.y) * rayUnitStepSize.y
	}

	tileFound := false
	maxDistance := 256.0
	distance := 0.0

	var side = 0
	for !tileFound && distance < maxDistance {

		if rayLength.x < rayLength.y {
			mapPos.x += Step.x
			distance += rayLength.x
			rayLength.x += rayUnitStepSize.x
			side = 0
		} else {
			mapPos.y += Step.y
			distance += rayLength.y
			rayLength.y += rayUnitStepSize.y
			side = 1
		}

		if mapPos.x > w.width || mapPos.x < 0 || mapPos.y > w.height || mapPos.y < 0 {
			break
		}

		t := w.getTile(mapPos.x, mapPos.y)
		if t != nil {
			if t.block {
				tileFound = true
			}
		}
	}

	perpWallDist := 256.0
	if tileFound {
		if side == 0 {
			perpWallDist = rayLength.x - rayUnitStepSize.x
		} else {
			perpWallDist = rayLength.y - rayUnitStepSize.y
		}
	}

	var wallx float64
	if side == 0 {
		wallx = rayStart.y + (perpWallDist * rayDir.y)
	} else {
		wallx = rayStart.x + (perpWallDist * rayDir.x)
	}
	wallx -= math.Floor(wallx)

	return ray{
		distance: perpWallDist,
		side:     side,
		wallx:    wallx,
		dir:      rayDir,
	}
}
