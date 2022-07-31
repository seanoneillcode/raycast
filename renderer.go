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
			LoadImage("floor-1.png"),
			LoadImage("ceiling.png"),
		},
	}
}

func (r *Renderer) Render(screen *ebiten.Image, w *World) {

	r.drawFloor(w)

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
			//if y < ScreenHeight/2 {
			//	r.SetActualPixel(float64(x), float64(y), skyColor)
			//} else {
			//	//r.SetActualPixel(float64(x), float64(y), grassColor)
			//}
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

func (r *Renderer) drawFloor(w *World) {
	for y := ScreenHeight / 2; y < ScreenHeight; y++ {
		// rayDir for leftmost ray (x = 0) and rightmost ray (x = w)
		rayDirX0 := w.playerDir.x - w.plane.x
		rayDirY0 := w.playerDir.y - w.plane.y
		rayDirX1 := w.playerDir.x + w.plane.x
		rayDirY1 := w.playerDir.y + w.plane.y

		// Current y position compared to the center of the screen (the horizon)
		p := y - ScreenHeight/2

		// Vertical position of the camera.
		// NOTE: with 0.5, it's exactly in the center between floor and ceiling,
		// matching also how the walls are being raycasted. For different values
		// than 0.5, a separate loop must be done for ceiling and floor since
		// they're no longer symmetrical.
		posZ := 0.5 * ScreenHeight

		// Horizontal distance from the camera to the floor for the current row.
		// 0.5 is the z position exactly in the middle between floor and ceiling.
		// NOTE: this is affine texture mapping, which is not perspective correct
		// except for perfectly horizontal and vertical surfaces like the floor.
		// NOTE: this formula is explained as follows: The camera ray goes through
		// the following two points: the camera itself, which is at a certain
		// height (posZ), and a point in front of the camera (through an imagined
		// vertical plane containing the screen pixels) with horizontal distance
		// 1 from the camera, and vertical position p lower than posZ (posZ - p). When going
		// through that point, the line has vertically traveled by p units and
		// horizontally by 1 unit. To hit the floor, it instead needs to travel by
		// posZ units. It will travel the same ratio horizontally. The ratio was
		// 1 / p for going through the camera plane, so to go posZ times farther
		// to reach the floor, we get that the total horizontal distance is posZ / p.
		rowDistance := posZ / float64(p)

		// calculate the real world step vector we have to add for each x (parallel to camera plane)
		// adding step by step avoids multiplications with a weight in the inner loop
		floorStepX := rowDistance * (rayDirX1 - rayDirX0) / ScreenWidth
		floorStepY := rowDistance * (rayDirY1 - rayDirY0) / ScreenWidth

		// real world coordinates of the leftmost column. This will be updated as we step to the right.
		floorX := w.playerPos.x + rowDistance*rayDirX0
		floorY := w.playerPos.y + rowDistance*rayDirY0

		for x := 0; x < ScreenWidth; x++ {
			// the cell coord is simply got from the integer parts of floorX and floorY
			cellX := (int)(floorX)
			cellY := (int)(floorY)

			// get the texture coordinate from the fractional part
			tx := (int)(TextureWidth*(floorX-float64(cellX))) & (TextureWidth - 1)
			ty := (int)(TextureHeight*(floorY-float64(cellY))) & (TextureHeight - 1)

			floorX += floorStepX
			floorY += floorStepY

			// floor
			img := r.textures[1]
			c := img.At(tx, ty)
			r.SetActualPixel(float64(x), float64(y), c)

			// ceiling
			img = r.textures[2]
			c = img.At(tx, ty)
			r.SetActualPixel(float64(x), float64(ScreenHeight-y-1), c)
		}

	}
}
