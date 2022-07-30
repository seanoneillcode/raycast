package raycast

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(imageFileName string) *ebiten.Image {
	//b, err := ioutil.ReadFile("res/" + imageFileName)
	//if err != nil {
	//	log.Fatalf("failed to open file: %v", err)
	//}
	file2, err := os.Open( imageFileName)
	if err != nil {
		fmt.Println("3 ", err)
		log.Fatal(err)
	}
	imageIn, _, err := image.Decode(file2)
	if err != nil {
		fmt.Println("4 ", err)
		log.Fatal(err)
	}
	//img, _, err := image.Decode(bytes.NewReader(b))
	//if err != nil {
	//	log.Fatal(err)
	//}
	return ebiten.NewImageFromImage(imageIn)
}
