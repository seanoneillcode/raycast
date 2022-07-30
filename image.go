package raycast

import (
	"image"
	"os"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func LoadImage(imageFileName string) *ebiten.Image {
	file, err := os.Open("res/" + imageFileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
