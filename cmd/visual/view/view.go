package view

import (
	"github.com/gdamore/tcell"
	"github.com/hinshun/astar/cmd/visual/terrain"
)

type Mode uint8

const (
	ModePlain Mode = iota
	ModeSwamp
	ModeWall
	ModeAgent
	ModeGoal
)

var (
	ModeToTile = map[Mode]terrain.Tile{
		ModePlain: terrain.TilePlain,
		ModeSwamp: terrain.TileSwamp,
		ModeWall:  terrain.TileWall,
		ModeAgent: terrain.TileAgent,
		ModeGoal:  terrain.TileGoal,
	}
)

type View struct {
	screen tcell.Screen
	mode   Mode
	grid   *terrain.Grid
}

func NewView(screen tcell.Screen, grid *terrain.Grid) *View {
	view := &View{
		screen: screen,
		mode:   ModeAgent,
		grid:   grid,
	}
	return view
}

func (v *View) Click(x, y int) {
	// Clicking out of bounds
	if v.OutOfBounds(x, y) {
		return
	}

	tile := ModeToTile[v.mode]

	wx := v.WorldX()
	wy := v.WorldY()
	v.grid.AddUpper(x-wx, y-wy, tile)
}

func (v *View) OutOfBounds(x, y int) bool {
	wx := v.WorldX()
	wy := v.WorldY()
	return x < wx || x >= wx+terrain.RoomSize || y < wy || y >= wy+terrain.RoomSize
}

func (v *View) SetMode(mode Mode) {
	v.mode = mode
}

func (v *View) WorldX() int {
	width, _ := v.screen.Size()
	mx := width / 2
	return mx - (terrain.RoomSize / 2)
}

func (v *View) WorldY() int {
	_, height := v.screen.Size()
	my := height / 2
	return my - (terrain.RoomSize / 2)
}

func (v *View) Update() {
	v.screen.Clear()

	wx := v.WorldX()
	wy := v.WorldY()

	merged := v.grid.Merged()
	for y, row := range merged {
		for x, tile := range row {
			style := terrain.TileToStyle[tile]
			ch := terrain.TileToRune[tile]

			switch tile {
			case terrain.TileWall:
				adjacent := false
				for _, position := range [][]int{
					{x, y - 1},
					{x + 1, y - 1},
					{x + 1, y},
					{x + 1, y + 1},
					{x, y + 1},
					{x - 1, y + 1},
					{x - 1, y},
					{x - 1, y - 1},
				} {
					if !v.OutOfBounds(wx+position[0], wy+position[1]) && merged[position[1]][position[0]] != terrain.TileWall {
						adjacent = true
					}
				}

				if !adjacent {
					continue
				}
			case terrain.TileAgent, terrain.TileGoal:
				lowerTile := v.grid.Lower[y][x]
				lowerStyle := terrain.TileToStyle[lowerTile]
				_, bg, _ := lowerStyle.Decompose()
				style = style.Background(bg)
			case terrain.TilePlain:
				if y == 0 {
					ch = '↑'
				} else if y == terrain.RoomSize-1 {
					ch = '↓'
				} else if x == 0 {
					ch = '←'
				} else if x == terrain.RoomSize-1 {
					ch = '→'
				}
			}

			v.screen.SetContent(wx+x, wy+y, ch, nil, style)
		}
	}

	v.screen.Show()
}

func (v *View) Resize(width, height int) {
}

func (v *View) Clear() {
	v.grid.Clear()
}
