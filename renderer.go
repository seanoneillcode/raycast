package raycast

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	image    *ebiten.Image
	textures []*ebiten.Image
}

func NewRenderer() *Renderer {
	return &Renderer{
		image: ebiten.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))),
		textures: []*ebiten.Image{
			LoadImage("stone.png"),
		},
	}
}

func (r *Renderer) Render(screen *ebiten.Image, w *World) {
	angleStep := FieldOfView / float64(NumRays)
	startAngle := w.playerDir - (angleStep * float64(NumRays/2))
	var rays []ray
	for rayIndex := 0; rayIndex < NumRays; rayIndex++ {
		angle := startAngle + (angleStep * float64(rayIndex))
		ra := getRay(w, angle)
		r.drawRay(ra, rayIndex)
	}
	r.renderDebug(w, rays)

	// final render to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GlobalScale, GlobalScale)
	screen.DrawImage(r.image, op)

	screen.DrawImage(r.textures[0], &ebiten.DrawImageOptions{})
}

func (r *Renderer) drawRay(ray ray, index int) {

	lineHeight := (int)(ScreenHeight / ray.distance)

	//calculate lowest and highest pixel to fill in current stripe
	drawStart := ScreenHeight/2 - lineHeight/2
	if drawStart < 0 {
		drawStart = 0
	}
	drawEnd := ScreenHeight/2 + lineHeight/2
	if drawEnd >= ScreenHeight {
		drawEnd = ScreenHeight - 1
	}

	width := ScreenWidth / NumRays

	for x := width * index; x < width*(index+1); x++ {
		for y := 0; y < ScreenHeight; y++ {
			if y > drawStart && y < drawEnd {
				if ray.side == 0 {
					r.SetActualPixel(float64(x), float64(y), halfWallColor)
				} else {
					r.SetActualPixel(float64(x), float64(y), wallColor)
				}
			} else {
				if y < ScreenHeight/2 {
					r.SetActualPixel(float64(x), float64(y), skyColor)
				} else {
					r.SetActualPixel(float64(x), float64(y), grassColor)
				}
			}
		}
	}
}

func (r *Renderer) SetPixel(x float64, y float64, c color.RGBA) {
	r.image.Set(int(x*TileSize), int(y*TileSize), c)
}

func (r *Renderer) SetActualPixel(x float64, y float64, c color.RGBA) {
	r.image.Set(int(x), int(y), c)
}

func (r *Renderer) renderDebug(w *World, rays []ray) {

	// draw tiles
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			t := w.getTile(x, y)
			if t.block {
				r.drawTile(x, y, blockColor)
			} else {
				r.drawEdgeTile(x, y)
			}
		}
	}

	// draw rays
	for _, ray := range rays {
		rayDir := point{
			x: math.Cos(ray.angle),
			y: math.Sin(ray.angle),
		}
		intersection := point{
			x: w.playerPos.x + (rayDir.x * ray.distance),
			y: w.playerPos.y + (rayDir.y * ray.distance),
		}
		r.drawHit(intersection)
	}

	// draw player
	r.SetPixel(w.playerPos.x, w.playerPos.y, playerColor)

	// draw player direction
	angle := w.playerDir * math.Pi / 3
	dirx := w.playerPos.x + (1 * math.Cos(angle))
	diry := w.playerPos.y + (1 * math.Sin(angle))
	r.SetPixel(dirx, diry, dirColor)
}

func (r *Renderer) drawHitCross(intersection point) {
	r.SetActualPixel((intersection.x*TileSize)-1, (intersection.y*TileSize)-1, crossColor)
	r.SetActualPixel((intersection.x*TileSize)+1, (intersection.y*TileSize)-1, crossColor)
	r.SetActualPixel((intersection.x*TileSize)+1, (intersection.y*TileSize)+1, crossColor)
	r.SetActualPixel((intersection.x*TileSize)-1, (intersection.y*TileSize)+1, crossColor)
}
func (r *Renderer) drawHit(intersection point) {
	r.SetActualPixel(intersection.x*TileSize, intersection.y*TileSize, crossColor)
}

func (r *Renderer) drawTile(tx int, ty int, c color.RGBA) {
	px := tx * TileSize
	py := ty * TileSize

	for x := px; x < (px + TileSize); x++ {
		for y := py; y < (py + TileSize); y++ {
			r.SetActualPixel(float64(x), float64(y), c)
		}
	}
}

func (r *Renderer) drawEdgeTile(tx int, ty int) {
	px := tx * TileSize
	py := ty * TileSize

	for x := px; x < (px + TileSize); x++ {
		for y := py; y < (py + TileSize); y++ {
			col := emptyColor
			if x == px {
				col = edgeColor
			}
			if y == py {
				col = edgeColor
			}
			r.SetActualPixel(float64(x), float64(y), col)
		}
	}
}
