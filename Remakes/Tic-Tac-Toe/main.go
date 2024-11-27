package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenwidth  = 600
	screenHeight = 600
	cellWidth    = 150
	cellHeight   = 150
)

var (
	game_font font.Face
	str_font  font.Face
)

type Cell struct {
	X, Y, W, H int
	char       string
	filled     bool
}

type Reset struct {
	X, Y, W, H int
	visible    bool
}

type Game struct {
	cells     [3][3]Cell
	clicks    int
	board     []string
	game_over bool
	reset     Reset
}

func init() {

	ttfData, err := os.ReadFile("assets/font.ttf")

	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(ttfData)

	if err != nil {
		log.Fatal(err)
	}

	game_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    60,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	str_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getWidth(str string, Font font.Face) int {

	width := 0

	for _, char := range str {

		w, _ := Font.GlyphAdvance(char)
		width += w.Ceil()
	}

	return width
}

func getCount(board []string) int {

	count := 0

	for _, char := range board {
		if char == "" {
			count++
		}
	}

	return count
}

func main() {

	ebiten.SetWindowSize(screenwidth, screenHeight)
	ebiten.SetWindowTitle("Tic-Tac-Toe - The Fourth")

	cells := [3][3]Cell{}

	board := []string{}

	for i := 0; i < 10; i++ {
		board = append(board, "")
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cells[i][j] = Cell{
				X:      75 + (cellWidth * j),
				Y:      75 + (cellHeight * i),
				W:      cellWidth,
				H:      cellHeight,
				char:   "",
				filled: false,
			}
		}
	}

	reset := Reset{
		X:       (screenwidth - 150) / 2,
		Y:       530,
		W:       150,
		H:       50,
		visible: false,
	}

	game := &Game{
		cells:     cells,
		clicks:    0,
		board:     board,
		game_over: false,
		reset:     reset,
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {

	g.isCellClicked()

	if g.CheckWinner("X") || g.CheckWinner("O") || g.CheckDraw() {
		g.game_over = true
		g.reset.visible = true

		for i, row := range g.cells {
			for j := range row {

				cell := &g.cells[i][j]

				cell.filled = true
			}
		}
	}

	g.isResetClicked()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, row := range g.cells {
		for _, cell := range row {

			vector.DrawFilledRect(screen, float32(cell.X), float32(cell.Y), float32(cell.W), float32(cell.H), color.White, false)
			vector.DrawFilledRect(screen, float32(cell.X+2), float32(cell.Y+2), float32(cell.W-4), float32(cell.H-4), color.Black, false)

			text.Draw(screen, cell.char, game_font, cell.X+55, cell.Y+90, color.White)
		}
	}

	x_str := "Player with X Plays next"
	o_str := "Player with O Plays next"

	x_win_str := "Player with X Wins!"
	o_win_str := "Player with O Wins!"
	draw_str := "It's a Draw!"

	if !(g.game_over) {

		width := getWidth(x_str, str_font)
		X := (screenwidth - width) / 2

		if g.clicks%2 == 0 {
			text.Draw(screen, x_str, str_font, X, 30, color.White)
		} else {
			text.Draw(screen, o_str, str_font, X, 30, color.White)
		}
	} else {

		win_width := getWidth(x_win_str, str_font)
		X := (screenwidth - win_width) / 2

		draw_width := getWidth(draw_str, str_font)
		X_D := (screenwidth - draw_width) / 2

		if g.CheckWinner("X") {

			text.Draw(screen, x_win_str, str_font, X, 30, color.White)

		} else if g.CheckWinner("O") {

			text.Draw(screen, o_win_str, str_font, X, 30, color.White)

		} else {

			text.Draw(screen, draw_str, str_font, X_D, 30, color.White)
		}
	}

	if g.reset.visible {

		vector.DrawFilledRect(screen, float32(g.reset.X), float32(g.reset.Y), float32(g.reset.W), float32(g.reset.H), color.White, false)
		vector.DrawFilledRect(screen, float32(g.reset.X+2), float32(g.reset.Y+2), float32(g.reset.W-4), float32(g.reset.H-4), color.Black, false)

		text.Draw(screen, "Reset", str_font, g.reset.X+45, g.reset.Y+40, color.White)
	}
}

func (g *Game) Layout(outsidewidth, outsideHeight int) (int, int) {

	return screenwidth, screenHeight
}

func (g *Game) UpdateBoard() {

	g.board[1] = g.cells[0][0].char
	g.board[2] = g.cells[0][1].char
	g.board[3] = g.cells[0][2].char
	g.board[4] = g.cells[1][0].char
	g.board[5] = g.cells[1][1].char
	g.board[6] = g.cells[1][2].char
	g.board[7] = g.cells[2][0].char
	g.board[8] = g.cells[2][1].char
	g.board[9] = g.cells[2][2].char

}

func (g *Game) CheckWinner(marker string) bool {

	if (g.board[1] == g.board[2] && g.board[2] == g.board[3] && g.board[1] == marker) ||
		(g.board[4] == g.board[5] && g.board[5] == g.board[6] && g.board[4] == marker) ||
		(g.board[7] == g.board[8] && g.board[8] == g.board[9] && g.board[9] == marker) ||
		(g.board[1] == g.board[4] && g.board[4] == g.board[7] && g.board[1] == marker) ||
		(g.board[2] == g.board[5] && g.board[5] == g.board[8] && g.board[2] == marker) ||
		(g.board[3] == g.board[6] && g.board[6] == g.board[9] && g.board[3] == marker) ||
		(g.board[1] == g.board[5] && g.board[5] == g.board[9] && g.board[1] == marker) ||
		(g.board[3] == g.board[5] && g.board[5] == g.board[7] && g.board[3] == marker) {

		return true
	}

	return false
}

func (g *Game) CheckDraw() bool {

	if !g.CheckWinner("X") && !g.CheckWinner("O") && getCount(g.board) <= 1 {
		return true
	}

	return false
}

func (g *Game) isCellClicked() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		x, y := ebiten.CursorPosition()

		for i, row := range g.cells {
			for j := range row {

				cell := &g.cells[i][j]

				if (x >= cell.X && x <= cell.X+cell.W) && (y >= cell.Y && y <= cell.Y+cell.H) {

					if !cell.filled {

						if g.clicks%2 == 0 {
							cell.char = "X"
						} else {
							cell.char = "O"
						}

						g.clicks++

						cell.filled = true

						g.UpdateBoard()
					}
				}
			}
		}
	}
}

func (g *Game) isResetClicked() {

	if g.reset.visible {

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

			x, y := ebiten.CursorPosition()

			if (x >= g.reset.X && x <= g.reset.X+g.reset.W) && (y >= g.reset.Y && y <= g.reset.Y+g.reset.H) {

				g.clicks = 0

				for i, row := range g.cells {
					for j := range row {

						cell := &g.cells[i][j]

						cell.filled = false
						cell.char = ""
					}
				}

				g.UpdateBoard()
				g.game_over = false
				g.reset.visible = false
			}
		}
	}
}
