package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth = 800
	ScreeHeight = 600
	CellSize    = 150
	ResetWidth  = 200
	ResetHeight = 100
)

var (
	cellFont font.Face
	strFont  font.Face
)

type Game struct {
	cells          [9]Cell
	clicks         int
	reset_visible  bool
	board          [10]string
	reset          Reset
	player_1_queue *node
	player_2_queue *node
	display_str    string
}

func init() {

	ttfData, err := os.ReadFile("font.ttf")

	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(ttfData)

	if err != nil {
		log.Fatal(err)
	}

	cellFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    75,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	strFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    50,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ebiten.SetWindowTitle("Infinite Tic-Tac-Toe")
	ebiten.SetWindowSize(ScreenWidth, ScreeHeight)

	cells := []Cell{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cells = append(cells, Cell{
				X:      (175 + (j * CellSize)),
				Y:      (75 + (i * CellSize)),
				W:      CellSize,
				H:      CellSize,
				marker: "",
				filled: false,
			})
		}
	}

	reset := Reset{
		Cell: Cell{
			X:      (ScreenWidth - ResetWidth)/2,
			Y:      530,
			W:      ResetWidth,
			H:      ResetHeight,
			marker: "Reset",
			filled: false,
		},
	}

	board := [10]string{}

	for i := 0; i < 10; i++ {
		board[i] = ""
	}

	game := &Game{
		cells:         [9]Cell(cells),
		clicks:        0,
		reset_visible: false,
		board:         board,
		reset:         reset,

		player_1_queue: &node{
			position: -1,
			next:     nil,
		},

		player_2_queue: &node{
			position: -1,
			next:     nil,
		},

		display_str: "It's X's Turn!",
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {

	g.isCellClicked()
	g.handleMarkerRemoval()
	g.UpdateStrings()

	if g.CheckWin("X") || g.CheckWin("O") {
		g.reset_visible = true

		for i := range g.cells {
			cell := &g.cells[i]

			cell.filled = true
		}
	}

	g.ResetClicked()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, cell := range g.cells {
		vector.DrawFilledRect(screen, float32(cell.X), float32(cell.Y), float32(CellSize), float32(CellSize), color.White, true)
		vector.DrawFilledRect(screen, float32(cell.X+2), float32(cell.Y+2), float32(CellSize-4), float32(CellSize-4), color.Black, true)

		width := getWidth(cell.marker, cellFont)
		text.Draw(screen, cell.marker, cellFont, (cell.X + (cell.W/2 - width/2)), (cell.Y + cell.H/2 + 10), color.White)
	}

	width := getWidth(g.display_str, strFont)
	text.Draw(screen, g.display_str, strFont, (ScreenWidth-width)/2, 50, color.White)

	if g.reset_visible {
		vector.DrawFilledRect(screen, float32(g.reset.X), float32(g.reset.Y), float32(g.reset.W), float32(g.reset.H), color.White, true)
		vector.DrawFilledRect(screen, float32(g.reset.X+1), float32(g.reset.Y+1), float32(g.reset.W-2), float32(g.reset.H-2), color.Black, true)

		width := getWidth(g.reset.marker, strFont)
		text.Draw(screen, g.reset.marker, strFont, g.reset.X+(g.reset.W/2-width/2), g.reset.Y+g.reset.H/2 + 10, color.White)
	}

}

func (g *Game) Layout(outsidewidth, outsideHeight int) (int, int) {

	return ScreenWidth, ScreeHeight
}
