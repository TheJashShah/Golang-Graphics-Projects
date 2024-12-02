package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 460
	screenHeight = 800
	cellSize     = 80
	buttonWidth  = 20
	buttonHeight = 40
)

var (
	str_font  font.Face
	btn_font  font.Face
	cell_font font.Face
)

type Cell struct {
	X, Y, W, H int
	str        string
	color      color.RGBA
}

type Button struct {
	X, Y, W, H int
	str        string
	color      color.RGBA
}

type Util struct {
	Button
}

type Game struct {
	cells   []Cell
	btns    []Button
	current int
	word    string
	guess   string
	utils   []Util
}

func init() {

	ttfData, err := os.ReadFile("assets/font.ttf")

	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(ttfData)

	if err != nil {
		log.Fatal(err)
	}

	cell_font, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    45,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	btn_font, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	str_font, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wordle - The Fifth")

	cells := []Cell{}

	for i := 0; i < 6; i++ {
		for j := 0; j < 5; j++ {
			cells = append(cells, Cell{
				X:     10 + (cellSize+10)*j,
				Y:     10 + (cellSize+10)*i,
				W:     cellSize,
				H:     cellSize,
				str:   "",
				color: color.RGBA{0, 0, 0, 255},
			})
		}
	}

	alphabets := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LIST := strings.Split(alphabets, "")

	btns := []Button{}

	for i := 0; i < 14; i++ {
		btns = append(btns, Button{
			X:     12 + (buttonWidth+12)*i,
			Y:     600,
			W:     buttonWidth,
			H:     buttonHeight,
			str:   LIST[i],
			color: color.RGBA{68, 65, 66, 1},
		})
	}

	for i := 0; i < 12; i++ {
		btns = append(btns, Button{
			X:     44 + (buttonWidth+12)*i,
			Y:     650,
			W:     buttonWidth,
			H:     buttonHeight,
			str:   LIST[14+i],
			color: color.RGBA{68, 65, 66, 1},
		})
	}

	utils := []Util{}

	utils = append(utils, Util{
		Button{
			X:     40,
			Y:     710,
			W:     120,
			H:     30,
			str:   "ENTER",
			color: color.RGBA{68, 65, 66, 1},
		},
	})

	utils = append(utils, Util{
		Button{
			X:     300,
			Y:     710,
			W:     120,
			H:     30,
			str:   "BACK",
			color: color.RGBA{68, 65, 66, 1},
		},
	})

	game := &Game{
		cells:   cells,
		btns:    btns,
		current: 0,
		word:    "SPEED",
		guess:   "",
		utils:   utils,
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {

	g.isLetterButtonPressed()
	g.areUtilsPressed()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, cell := range g.cells {

		vector.DrawFilledRect(screen, float32(cell.X), float32(cell.Y), float32(cell.W), float32(cell.H), color.White, false)
		vector.DrawFilledRect(screen, float32(cell.X+1), float32(cell.Y+1), float32(cell.W-2), float32(cell.H-2), cell.color, false)

		text.Draw(screen, cell.str, cell_font, cell.X+25, cell.Y+50, color.White)
	}

	for _, btn := range g.btns {

		vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), btn.color, false)

		text.Draw(screen, btn.str, btn_font, btn.X+4, btn.Y+20, color.White)
	}

	for _, btn := range g.utils {

		vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), btn.color, false)

		text.Draw(screen, btn.str, str_font, btn.X+25, btn.Y+25, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) isLetterButtonPressed() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		x, y := ebiten.CursorPosition()

		for i := range g.btns {

			btn := &g.btns[i]

			if (x >= btn.X && x <= btn.X+btn.W) && (y >= btn.Y && y <= btn.Y+btn.H) {

				if g.current <= 29 {
					g.cells[g.current].str = btn.str
				}

				if g.current < 30 {
					g.current++
				}
			}
		}
	}
}

func (g *Game) GetUserWord() {

	if g.current%5 == 0 && g.current != 0 {

		g.guess = ""

		for i := g.current - 5; i < g.current; i++ {
			g.guess += g.cells[i].str
		}

		fmt.Println(g.guess)

	} else {
		return
	}
}

func (g *Game) Check() {

	if g.word == g.guess {

		for i := g.current - 5; i < g.current; i++ {
			g.cells[i].color = color.RGBA{22, 248, 22, 255}
		}

		for i := range g.btns {
			btn := &g.btns[i]

			if strings.Contains(g.word, btn.str) {
				btn.color = color.RGBA{22, 248, 22, 255}
			}
		}
		fmt.Println("You have Guessed the Word!")
		return
	}

	usedPositions := make([]bool, len(g.word))

	for i := 0; i < 5; i++ {
		if g.guess[i] == g.word[i] {
			g.cells[i+g.current-5].color = color.RGBA{22, 248, 22, 255}
			usedPositions[i] = true
		}
	}

	// for i := 0; i < 5; i++ {
	// 	if g.guess[i] != g.word[i] {
	// 		for j := 0; j < 5; j++ {
	// 			if g.guess[i] == g.word[j] && !usedPositions[j] {
	// 				g.cells[i+g.current-5].color = color.RGBA{243, 245, 39, 255}
	// 				usedPositions[j] = true
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	for i := 0; i < 5; i++ {
		if g.guess[i] != g.word[i] && strings.Contains(g.word, string(g.guess[i])) {
			g.cells[g.current+i-5].color = color.RGBA{243, 245, 39, 255}
		}
	}

	for a := range g.btns {
		btn := &g.btns[a]
		for i := 0; i < 5; i++ {
			if btn.str == string(g.guess[i]) {
				if g.guess[i] == g.word[i] {
					btn.color = color.RGBA{22, 248, 22, 255}
				} else if strings.Contains(g.word, btn.str) {
					btn.color = color.RGBA{243, 245, 39, 255}
				}
			}
		}
	}
}

func (g *Game) areUtilsPressed() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		x, y := ebiten.CursorPosition()

		for i := range g.utils {

			btn := &g.utils[i]

			if (x >= btn.X && x <= btn.X+btn.X+btn.W) && (y >= btn.Y && y <= btn.Y+btn.H) {

				if i == 0 {

					g.GetUserWord()
					g.Check()

				} else if i == 1 {
					if g.current > 0 {
						g.current--
						g.cells[g.current].str = ""
					}
				}
			}
		}
	}
}
// FINISHING TOUCHES ARE LEFT BUT I'M BORED.
// THINGS REMAINING(IF ANYONE WANTS TO COMPLETE):
// 1. RANDOM WORD GENERATION
// 2. GAME OVER WHEN 6 TRIES ARE COMPLETED.
// 3. RESET OPTION