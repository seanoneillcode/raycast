package raycast

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowWidth   = 800
	WindowHeight  = 800
	ScreenWidth   = 256
	ScreenHeight  = 256
	PlayerWidth   = 4
	GlobalScale   = 1
	MapSize       = 16
	NumRays       = 256
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
		world:            NewWorld(MapSize, MapSize),
		renderer:         NewRenderer(),
		lastUpdateCalled: time.Now(),
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
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
