package raycast

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"math"
	"sort"
)

type Renderer struct {
	background      *ebiten.Image
	image           *ebiten.Image
	weaponAnimation *animation
	textures        map[string]image.Image
	zbuffer         []float64
	commonFont      font.Face
}

var fakeLightEnabled = false
var screenFlashEnabled = true

func NewRenderer() *Renderer {
	return &Renderer{
		background: ebiten.NewImageFromImage(LoadImage("background.png")),
		image:      ebiten.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))),
		textures: map[string]image.Image{
			"wall":              LoadImage("wall-2.png"),
			"floor":             LoadImage("floor-1.png"),
			"ceiling":           LoadImage("ceiling.png"),
			"eye":               LoadImage("sprite.png"),
			"bullet":            LoadImage("bullet.png"),
			"door":              LoadImage("door.png"),
			"door-floor":        LoadImage("door-floor.png"),
			"door-wall":         LoadImage("door-wall.png"),
			"soul":              LoadImage("soul.png"),
			"ammo":              LoadImage("ammo.png"),
			"ammo-icon":         LoadImage("ammo-icon.png"),
			"health":            LoadImage("health.png"),
			"health-icon":       LoadImage("health-icon.png"),
			"bullet-hit":        LoadImage("bullet-hit.png"),
			"portal":            LoadImage("portal.png"),
			"weapon-staff":      LoadImage("weapon-staff.png"),
			"enemy-ball-move":   LoadImage("enemy-ball-move.png"),
			"enemy-ball-hurt":   LoadImage("enemy-ball-hurt.png"),
			"enemy-ball-attack": LoadImage("enemy-ball-attack.png"),
			"enemy-ball-die":    LoadImage("enemy-ball-die.png"),
			"enemy-blue-move":   LoadImage("enemy-blue-move.png"),
			"enemy-blue-hurt":   LoadImage("enemy-blue-hurt.png"),
			"enemy-blue-attack": LoadImage("enemy-blue-attack.png"),
			"enemy-blue-die":    LoadImage("enemy-blue-die.png"),
			"candlestick":       LoadImage("candlestick.png"),
		},
		zbuffer: make([]float64, ScreenWidth),
	}
}

func (r *Renderer) Render(screen *ebiten.Image, w *World) {
	r.image.Clear()
	r.drawSky(w)

	r.drawFloorAndCeiling(w)

	for rayIndex := 0; rayIndex < NumRays; rayIndex++ {
		// cameraX goes from -1 to +1 (very roughly)
		cameraX := 2*(float64(rayIndex)/float64(NumRays)) - 1
		ra := calculateRay(w, cameraX)
		r.drawRay(ra, rayIndex)
		r.zbuffer[rayIndex] = ra.distance
	}

	r.drawSprites(w)

	r.drawHud(w)
	r.drawWeapon(w)
	r.drawMiniMap(w)

	// final render to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GlobalScale, GlobalScale)
	if screenFlashEnabled && w.player.screenFlashTimer > 0 {
		screen.Fill(w.player.screenFlashColor)
		scale := 1 - ((w.player.screenFlashTimer / screenFlashTime) / 1)
		op.ColorM.Scale(1, 1, 1, scale)
	}
	screen.DrawImage(r.image, op)
}

func (r *Renderer) drawSky(w *World) {
	angle := math.Atan2(w.player.dir.y, w.player.dir.x)
	angle = (angle + (math.Pi)) / (2 * math.Pi)

	var doubleWidth = ScreenWidth * 2

	for x := 0; x < ScreenWidth; x++ {
		for y := 0; y < ScreenHeight; y++ {

			xoffset := x + int(8*angle*ScreenWidth)
			xoffset = xoffset % doubleWidth
			if xoffset > doubleWidth {
				xoffset -= doubleWidth
			}
			if xoffset < 0 {
				xoffset += doubleWidth
			}
			c := r.background.At(xoffset, y)
			r.SetPixel(float64(x), float64(y), c)
		}
	}
}

func (r *Renderer) drawWeapon(w *World) {
	pos := vector{
		x: 192 - 64,
		y: 192 - 64,
	}
	img := r.textures["weapon-staff"]
	r.drawScaledImage(pos, img, w.player.weaponAnimation.currentFrame, TextureWidth*2, 2)
}

func (r *Renderer) drawHud(w *World) {
	ammoPos := vector{
		x: 24,
		y: 8,
	}
	ammoIcon := r.textures["ammo-icon"]
	r.drawScaledImage(ammoPos, ammoIcon, 0, TextureWidth, 1)

	healthPos := vector{
		x: 256 - 32 - 8,
		y: 8,
	}
	healthIcon := r.textures["health-icon"]
	r.drawScaledImage(healthPos, healthIcon, 0, TextureWidth, 1)

	RenderText(r.image, fmt.Sprintf("%d", w.player.ammo), int(ammoPos.x+8), 7)
	RenderText(r.image, fmt.Sprintf("%d", w.player.health), int(healthPos.x+8), 7)

	//RenderText(r.image, "find the portal to escape the maze!\nlots of love,\nbad wizard.", 32, 32)
}

func (r *Renderer) drawScaledImage(pos vector, img image.Image, frame int, textureWidth int, scale int) {
	frameOffsetX := frame * textureWidth
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := img.At(x+frameOffsetX, y)
			for q := x * scale; q < (x*scale)+scale; q++ {
				for z := y * scale; z < (y*scale)+scale; z++ {
					r.SetPixel(float64(q)+pos.x, float64(z)+pos.y, c)
				}
			}
		}
	}
}

func (r *Renderer) drawSprites(w *World) {

	var sprites []*sprite

	for _, e := range w.enemies {
		sprites = append(sprites, e.entity.CurrentSprite())
	}
	for _, b := range w.bullets {
		sprites = append(sprites, b.entity.CurrentSprite())
	}
	for _, b := range w.pickups {
		sprites = append(sprites, b.entity.CurrentSprite())
	}
	for _, b := range w.effects {
		sprites = append(sprites, b.entity.CurrentSprite())
	}
	for _, b := range w.portals {
		sprites = append(sprites, b.entity.CurrentSprite())
	}
	for _, b := range w.scenery {
		sprites = append(sprites, b.entity.CurrentSprite())
	}

	for _, s := range sprites {
		s.distance = (w.player.pos.x-s.pos.x)*(w.player.pos.x-s.pos.x) + (w.player.pos.y-s.pos.y)*(w.player.pos.y-s.pos.y)
	}

	sort.Slice(sprites, func(i, j int) bool {
		return sprites[i].distance > sprites[j].distance
	})

	for _, s := range sprites {
		spriteX := s.pos.x - w.player.pos.x
		spriteY := s.pos.y - w.player.pos.y

		//transform sprite with the inverse camera matrix
		// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
		// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
		// [ planeY   dirY ]                                          [ -planeY  planeX ]

		invDet := 1.0 / (w.player.plane.x*w.player.dir.y - w.player.dir.x*w.player.plane.y) //required for correct matrix multiplication

		transformX := invDet * (w.player.dir.y*spriteX - w.player.dir.x*spriteY)
		transformY := invDet * (-w.player.plane.y*spriteX + w.player.plane.x*spriteY) //this is actually the depth inside the screen, that what Z is in 3D, the distance of sprite to player, matching sqrt(spriteDistance[i])

		spriteScreenX := int((NumRays / 2) * (1 + transformX/transformY))

		//parameters for scaling and moving the sprites
		var uDiv = 1.0
		var vDiv = 1.0
		var vMove = s.height * TextureHeight
		vMoveScreen := int(vMove / transformY)

		//calculate height of the sprite on screen
		spriteHeight := int(math.Abs(ScreenHeight/(transformY)) / vDiv) //using "transformY" instead of the real distance prevents fisheye
		//calculate lowest and highest pixel to fill in current stripe
		drawStartY := (-spriteHeight/2 + ScreenHeight/2) + vMoveScreen
		if drawStartY < 0 {
			drawStartY = 0
		}
		drawEndY := (spriteHeight/2 + ScreenHeight/2) + vMoveScreen
		if drawEndY >= ScreenHeight {
			drawEndY = ScreenHeight - 1
		}

		//calculate width of the sprite
		spriteWidth := int(math.Abs(ScreenHeight/(transformY)) / uDiv) // same as height of sprite, given that it's square
		drawStartX := -spriteWidth/2 + spriteScreenX
		if drawStartX < 0 {
			drawStartX = 0
		}
		drawEndX := spriteWidth/2 + spriteScreenX
		if drawEndX > NumRays {
			drawEndX = NumRays
		}

		//loop through every vertical stripe of the sprite on screen
		for stripe := drawStartX; stripe < drawEndX; stripe++ {
			texX := int(256*(stripe-(-spriteWidth/2+spriteScreenX))*TextureWidth/spriteWidth) / 256
			//the conditions in the if are:
			//1) it's in front of camera plane so you don't see things behind you
			//2) ZBuffer, with perpendicular distance
			if transformY > 0 && transformY < r.zbuffer[stripe] {
				for y := drawStartY; y < drawEndY; y++ { //for every pixel of the current stripe
					d := (y-vMoveScreen)*256 - ScreenHeight*128 + spriteHeight*128 //256 and 128 factors to avoid floats
					texY := ((d * TextureHeight) / spriteHeight) / 256

					img := r.textures[s.image]
					frameOffsetX := 0
					if s.animation != nil {
						frameOffsetX = s.animation.currentFrame * TextureWidth
					}
					c := img.At(texX+frameOffsetX, texY)
					r.SetPixel(float64(stripe), float64(y), c)
				}
			}
		}

	}
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

	var texX = int(ray.wallX * TextureWidth)
	//flip textures if looking in opposite direction
	if ray.side == 0 && ray.dir.x > 0 {
		texX = TextureWidth - texX - 1
	}
	if ray.side == 1 && ray.dir.y < 0 {
		texX = TextureWidth - texX - 1
	}
	texture := "wall"
	if ray.texture != "" {
		texture = ray.texture
	}
	img := r.textures[texture]
	if img == nil {
		fmt.Println("sadfa")
	}

	x := index
	step := float64(TextureHeight) / float64(lineHeight)
	texPos := float64(drawStart-ScreenHeight/2+lineHeight/2) * step

	for y := drawStart; y < drawEnd; y++ {
		texY := int(texPos) & (TextureHeight - 1)
		texPos += step

		c := img.At(texX, texY)

		rgba := color.RGBAModel.Convert(c).(color.RGBA)
		rgba = fakeLight(rgba, ray.distance)
		if ray.side == 0 {
			rgba.R = rgba.R / 2
			rgba.G = rgba.G / 2
			rgba.B = rgba.B / 2
		}
		c = rgba
		r.SetPixel(float64(x), float64(y), c)
	}
}

func fakeLight(c color.RGBA, distance float64) color.RGBA {
	if !fakeLightEnabled {
		return c
	}

	const maxLightDistance = 10.0
	value := maxLightDistance - distance
	if value < 0 {
		value = 0
	}
	percentVal := value / maxLightDistance

	newR := ((float64(c.R) / 255.0) * percentVal) * (255)
	newG := ((float64(c.G) / 255.0) * percentVal) * (255)
	newB := ((float64(c.B) / 255.0) * percentVal) * (255)

	c.R = uint8(newR)
	c.G = uint8(newG)
	c.B = uint8(newB)

	return c
}

func (r *Renderer) SetPixel(x float64, y float64, c color.Color) {
	_, _, _, a := c.RGBA()
	if a == 0 {
		return
	}
	r.image.Set(int(x), int(y), c)
}

func (r *Renderer) drawFloorAndCeiling(w *World) {
	for y := ScreenHeight / 2; y < ScreenHeight; y++ {
		// rayDir for leftmost ray (x = 0) and rightmost ray (x = w)
		rayDirX0 := w.player.dir.x - w.player.plane.x
		rayDirY0 := w.player.dir.y - w.player.plane.y
		rayDirX1 := w.player.dir.x + w.player.plane.x
		rayDirY1 := w.player.dir.y + w.player.plane.y

		// Current y position compared to the center of the screen (the horizon)
		p := y - ScreenHeight/2 + 1

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
		floorX := w.player.pos.x + rowDistance*rayDirX0
		floorY := w.player.pos.y + rowDistance*rayDirY0

		for x := 0; x < ScreenWidth; x++ {
			// the cell coord is simply got from the integer parts of floorX and floorY
			cellX := (int)(floorX)
			cellY := (int)(floorY)

			t := w.getTile(cellX, cellY)
			floorTex := ""
			ceilingTex := ""
			if t != nil {
				floorTex = t.floorTex
				ceilingTex = t.ceilingTex
			}

			// get the texture coordinate from the fractional part
			tx := (int)(TextureWidth*(floorX-float64(cellX))) & (TextureWidth - 1)
			ty := (int)(TextureHeight*(floorY-float64(cellY))) & (TextureHeight - 1)

			floorX += floorStepX
			floorY += floorStepY

			if floorTex != "" {
				img := r.textures[floorTex]
				c := img.At(tx, ty)

				rgba := color.RGBAModel.Convert(c).(color.RGBA)
				rgba = fakeLight(rgba, rowDistance)
				r.SetPixel(float64(x), float64(y), rgba)
			}
			if ceilingTex != "" {
				img := r.textures[ceilingTex]
				c := img.At(tx, ty)
				rgba := color.RGBAModel.Convert(c).(color.RGBA)
				rgba = fakeLight(rgba, rowDistance)
				r.SetPixel(float64(x), float64(ScreenHeight-y-1), rgba)
			}

		}

	}
}

func (r *Renderer) drawMiniMap(w *World) {

	if !w.player.showMiniMap {
		return
	}

	const miniWidth = 32
	const halfMiniWidth = miniWidth / 2

	mapPlayerPos := mapPos{
		x: int(w.player.pos.x),
		y: int(w.player.pos.y),
	}

	screenPos := vector{
		x: 8,
		y: 216,
	}

	for tx := 0; tx < miniWidth; tx += 1 {
		for ty := 0; ty < miniWidth; ty += 1 {
			x := mapPlayerPos.x + tx - halfMiniWidth
			y := mapPlayerPos.y + ty - halfMiniWidth

			t := w.getTile(x, y)
			if t == nil {
				continue
			}
			if !t.seen {
				continue
			}

			c := emptyColor
			if t.block {
				c = blockColor
			}
			if t.door {
				c = doorColor
			}

			r.drawChunkyPixel(screenPos.x+float64(x*2), screenPos.y+float64(y*2), c)
			if tx == halfMiniWidth && ty == halfMiniWidth {
				r.drawChunkyPixel(screenPos.x+float64(x*2), screenPos.y+float64(y*2), playerColor)
			}
		}
	}
}

func (r *Renderer) drawChunkyPixel(x float64, y float64, c color.RGBA) {
	r.SetPixel(x, y, c)
	r.SetPixel(x+1, y, c)
	r.SetPixel(x, y+1, c)
	r.SetPixel(x+1, y+1, c)
}
