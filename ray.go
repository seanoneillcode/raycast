package raycast

import "math"

type ray struct {
	side     int
	distance float64
	angle    float64
}

func getRay(w *World, angle float64) ray {
	rayStart := point{
		x: w.playerPos.x,
		y: w.playerPos.y,
	}

	rayDir := point{
		x: math.Cos(angle),
		y: math.Sin(angle),
	}

	rayUnitStepSize := point{
		x: math.Sqrt(1 + (rayDir.y/rayDir.x)*(rayDir.y/rayDir.x)),
		y: math.Sqrt(1 + (rayDir.x/rayDir.y)*(rayDir.x/rayDir.y)),
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

	modDistance := 256.0
	if tileFound {
		modDistance = rayLength.x - rayUnitStepSize.x
		if side == 1 {
			modDistance = rayLength.y - rayUnitStepSize.y
		}
	}

	return ray{
		distance: modDistance * math.Cos(angle-w.playerDir),
		side:     side,
		angle:    angle,
	}
}
