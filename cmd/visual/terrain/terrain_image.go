package terrain

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/pkg/errors"
)

const (
	cloudfrontURL = "https://d3os7yery2usni.cloudfront.net/map"
)

var (
	rgbaToTile = map[rgba]Tile{
		rgba{0x2b2b, 0x2b2b, 0x2b2b, 0xffff}: TilePlain,
		rgba{0x2323, 0x2525, 0x1313, 0xffff}: TileSwamp,
		rgba{0x3232, 0x3232, 0x3232, 0xffff}: TilePlain, // TileExit
		rgba{0, 0, 0, 0xffff}:                TileWall,
	}
)

type rgba struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

func TerrainImageToTerrain(shard, room string) ([][]Tile, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s.png", cloudfrontURL, shard, room))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get terrain image")
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode terrain image")
	}

	terrain := make([][]Tile, RoomSize)
	for i := range terrain {
		terrain[i] = make([]Tile, RoomSize)
	}

	for i := 0; i < RoomSize; i++ {
		for j := 0; j < RoomSize; j++ {
			// The terrain image is 150x150, so we skip to every third pixel.
			r, g, b, a := img.At(i*3, j*3).RGBA()
			terrain[j][i] = rgbaToTile[rgba{r, g, b, a}]
		}
	}

	return terrain, nil
}
