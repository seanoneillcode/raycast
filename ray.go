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
	for !tileFound && distance < maxDistance {

		if rayLength.x < rayLength.y {
			rayMapPos.x += step.x
			distance += rayLength.x
			rayLength.x += rayUnitStepSize.x
			side = 0
		} else {
			rayMapPos.y += step.y
			distance += rayLength.y
			rayLength.y += rayUnitStepSize.y
			side = 1
		}

		t := w.getTile(rayMapPos.x, rayMapPos.y)
		if t != nil {
			texture = t.wallTex
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
