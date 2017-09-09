package terrain

import (
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
)

type Tile uint8

const (
	RoomSize = 50

	TileNone Tile = iota
	TilePlain
	TileSwamp
	TileWall
	TileAgent
	TileGoal
)

var (
	ColorBase = tcell.NewRGBColor(0, 0, 95)
	ColorBlue = tcell.NewRGBColor(122, 158, 192)
	ColorTan  = tcell.NewRGBColor(88, 88, 0)

	TileToRune = map[Tile]rune{
		TilePlain: '.',
		TileSwamp: '~',
		TileWall:  '#',
		TileAgent: '@',
		TileGoal:  '!',
	}

	TileToStyle = map[Tile]tcell.Style{
		TilePlain: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(ColorBase),
		TileSwamp: tcell.StyleDefault.Foreground(ColorBlue).Background(ColorTan),
		TileWall:  tcell.StyleDefault.Foreground(tcell.ColorGrey).Background(tcell.ColorWhite),
		TileAgent: tcell.StyleDefault.Foreground(tcell.ColorWhite),
		TileGoal:  tcell.StyleDefault.Foreground(tcell.ColorRed),
	}
)

type Grid struct {
	Lower [][]Tile
	Upper [][]Tile
}

func NewGrid(shard, room string) (*Grid, error) {
	terrain, err := TerrainImageToTerrain(shard, room)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get terrain")
	}

	grid := &Grid{Lower: terrain}
	grid.Clear()

	return grid, nil
}

func (g *Grid) AddUpper(x, y int, tile Tile) {
	g.Upper[y][x] = tile
}

func (g *Grid) Merged() [][]Tile {
	merged := make([][]Tile, RoomSize)
	for y := range g.Lower {
		merged[y] = make([]Tile, RoomSize)
		for x := range g.Lower[y] {
			if g.Upper[y][x] != TileNone {
				merged[y][x] = g.Upper[y][x]
			} else {
				merged[y][x] = g.Lower[y][x]
			}
		}
	}

	return merged
}

func (g *Grid) Clear() {
	upper := make([][]Tile, RoomSize)
	for y := range upper {
		upper[y] = make([]Tile, RoomSize)
		for x := range upper[y] {
			upper[y][x] = TileNone
		}
	}
	g.Upper = upper
}
