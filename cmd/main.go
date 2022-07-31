package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"raycast.com"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	g := raycast.NewGame()

	ebiten.SetWindowSize(raycast.WindowWidth, raycast.WindowHeight)
	ebiten.SetWindowTitle("Raycast DEMO")
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
