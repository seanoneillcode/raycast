package raycast

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type Renderer struct {
	image    *ebiten.Image
	textures []image.Image
}

func NewRenderer() *Renderer {
	return &Renderer{
		image: ebiten.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))),
		textures: []image.Image{
			LoadImage("wall-2.png"),
			LoadImage("stone.png"),
		},
	}
}

func (r *Renderer) Render(screen *ebiten.Image, w *World) {
	for rayIndex := 0; rayIndex < NumRays; rayIndex++ {
		// cameraX goes from -1 to +1 (very roughly)
		cameraX := 2*(float64(rayIndex)/float64(NumRays)) - 1
		ra := calculateRay(w, cameraX)
		r.drawRay(ra, rayIndex)
	}

	// final render to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GlobalScale, GlobalScale)
	screen.DrawImage(r.image, op)
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

	var texX = int(ray.wallx * TextureWidth)
	//flip textures if looking in opposite direction
	if ray.side == 0 && ray.dir.x > 0 {
		texX = TextureWidth - texX - 1
	}
	if ray.side == 1 && ray.dir.y < 0 {
		texX = TextureWidth - texX - 1
	}

	x := index
	step := float64(TextureHeight) / float64(lineHeight)
	texPos := float64(drawStart-ScreenHeight/2+lineHeight/2) * step

	for y := 0; y < ScreenHeight; y++ {
		if y > drawStart && y < drawEnd {

			texY := int(texPos) & (TextureHeight - 1)
			texPos += step
			img := r.textures[0]
			c := img.At(texX, texY)
			if ray.side == 0 {
				rgba := color.RGBAModel.Convert(c).(color.RGBA)
				rgba.R = rgba.R / 2
				rgba.G = rgba.G / 2
				rgba.B = rgba.B / 2
				c = rgba
			}
			r.SetActualPixel(float64(x), float64(y), c)
		} else {
			if y < ScreenHeight/2 {
				r.SetActualPixel(float64(x), float64(y), skyColor)
			} else {
				r.SetActualPixel(float64(x), float64(y), grassColor)
			}
		}
	}
}

func (r *Renderer) SetPixel(x float64, y float64, c color.RGBA) {
	r.image.Set(int(x*TileSize), int(y*TileSize), c)
}

func (r *Renderer) SetActualPixel(x float64, y float64, c color.Color) {
	r.image.Set(int(x), int(y), c)
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
