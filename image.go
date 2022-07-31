package raycast

import (
	"image"
	"os"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func LoadEbitenImage(imageFileName string) *ebiten.Image {
	return ebiten.NewImageFromImage(LoadImage(imageFileName))
}

func LoadImage(imageFileName string) image.Image {
	file, err := os.Open("res/" + imageFileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
