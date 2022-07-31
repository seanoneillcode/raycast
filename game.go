package raycast

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WindowWidth   = 1024
	WindowHeight  = 768
	ScreenWidth   = 256
	ScreenHeight  = 192
	PlayerWidth   = 4
	GlobalScale   = 1
	MapSize       = 8
	TileSize      = 1
	NumRays       = 256
	FieldOfView   = (70) * (math.Pi / 180)
	TextureWidth  = 32
	TextureHeight = 32
)

type Game struct {
	world            *World
	renderer         *Renderer
	lastUpdateCalled time.Time
}

func NewGame() *Game {
	return &Game{
		world:    NewWorld(MapSize, MapSize),
		renderer: NewRenderer(),
	}
}

func (g *Game) Update() error {
	delta := time.Now().Sub(g.lastUpdateCalled).Milliseconds()
	g.lastUpdateCalled = time.Now()
	err := g.world.Update(float64(delta))
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.world)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
