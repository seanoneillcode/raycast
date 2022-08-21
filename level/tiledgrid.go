package level

import (
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

const (
	resourceDirectory = "res/maps/"
)

type TiledGrid struct {
	Layers            []*Layer            `json:"layers"`
	TileSetReferences []*TileSetReference `json:"tilesets"`
	TileSet           []*TileSet
}

type Layer struct {
	Data    []int         `json:"Data"`
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Objects []TiledObject `json:"objects"`
}

type TiledObject struct {
	Name       string            `json:"Name"`
	Type       string            `json:"type"`
	X          int               `json:"x"`
	Y          int               `json:"y"`
	Properties []*TileConfigProp `json:"properties"`
}

type TileSetReference struct {
	Source   string `json:"source"`
	FirstGid int    `json:"firstgid"`
}

type TileSet struct {
	ImageFileName string `json:"image"`
	ImageWidth    int    `json:"imagewidth"`
	ImageHeight   int    `json:"imageheight"`
	numTilesX     int
	numTilesY     int
	FirstGid      int
	Tiles         []*TileConfig `json:"tiles"`
	image         *ebiten.Image
}

type TileConfig struct {
	Id         int               `json:"id"`
	Properties []*TileConfigProp `json:"properties"`
}

type TileConfigProp struct {
	Name  string      `json:"Name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func NewTileGrid(fileName string) *TiledGrid {
	var tiledGrid TiledGrid

	configFile, err := os.Open(filepath.Join(resourceDirectory, fileName))
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&tiledGrid); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	tiledGrid.TileSet = []*TileSet{}
	for _, ref := range tiledGrid.TileSetReferences {
		tiledGrid.TileSet = append(tiledGrid.TileSet, loadTileSet(ref))
	}

	return &tiledGrid
}

func loadTileSet(ref *TileSetReference) *TileSet {
	tileSetConfigFile, err := os.Open(filepath.Join(resourceDirectory, ref.Source))
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	var tileSet TileSet
	jsonParser := json.NewDecoder(tileSetConfigFile)
	if err = jsonParser.Decode(&tileSet); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	tileSet.FirstGid = ref.FirstGid
	return &tileSet
}

func (tg *TiledGrid) getTileSetForIndex(index int) *TileSet {
	for i, tileSet := range tg.TileSet {
		if i == len(tg.TileSet)-1 || tg.TileSet[i+1].FirstGid > index {
			return tileSet
		}
	}
	// should never happen
	return nil
}

type ObjectData struct {
	Name       string
	ObjectType string
	X          int
	Y          int
	Properties []*ObjectProperty
}

type ObjectProperty struct {
	Name    string
	ObjType string
	Value   interface{}
}

func (tg *TiledGrid) GetObjectData() []*ObjectData {
	var ods []*ObjectData
	for _, l := range tg.Layers {
		for _, obj := range l.Objects {
			od := &ObjectData{
				Name:       obj.Name,
				ObjectType: obj.Type,
				X:          obj.X,
				Y:          obj.Y,
				Properties: []*ObjectProperty{},
			}
			for _, p := range obj.Properties {
				od.Properties = append(od.Properties, &ObjectProperty{
					Name:    p.Name,
					ObjType: p.Type,
					Value:   p.Value,
				})
			}
			ods = append(ods, od)
		}
	}
	return ods
}

type TileData struct {
	X          int
	Y          int
	Block      bool
	Door       bool
	North      bool
	WallTex    string
	FloorTex   string
	DoorTex    string
	CeilingTex string
}

func (tg *TiledGrid) GetTileData(x int, y int) *TileData {
	td := TileData{
		X: x,
		Y: y,
	}
	index := (y * tg.Layers[0].Width) + x

	if index < 0 || index >= len(tg.Layers[0].Data) {
		// no tile here
		return &td
	}

	if x < 0 || y < 0 {
		return &td
	}

	tileSetIndex := tg.Layers[0].Data[index]
	// I think this means there's nothing there ???
	if tileSetIndex == 0 {
		return &td
	}

	ts := tg.getTileSetForIndex(tileSetIndex)

	for _, tile := range ts.Tiles {
		if tile.Id == tileSetIndex-ts.FirstGid {
			for _, prop := range tile.Properties {
				if prop.Name == "block" && prop.Value != nil {
					td.Block = (prop.Value).(bool)
				}
				if prop.Name == "door" && prop.Value != nil {
					td.Door = (prop.Value).(bool)
				}
				if prop.Name == "north" && prop.Value != nil {
					td.North = (prop.Value).(bool)
				}
				if prop.Name == "wallTex" && prop.Value != nil {
					td.WallTex = (prop.Value).(string)
				}
				if prop.Name == "floorTex" && prop.Value != nil {
					td.FloorTex = (prop.Value).(string)
				}
				if prop.Name == "ceilingTex" && prop.Value != nil {
					td.CeilingTex = (prop.Value).(string)
				}
				if prop.Name == "doorTex" && prop.Value != nil {
					td.DoorTex = (prop.Value).(string)
				}
			}
			break
		}
	}

	return &td
}
