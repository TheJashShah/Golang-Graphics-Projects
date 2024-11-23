package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	game_font font.Face
)

const (
	screenWidth  = 640
	screenHeight = 480
	buttonWidth  = 200
	buttonHeight = 40
)

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
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func generateRandom() int {

	return rand.Intn(9)
}

type Button struct {
	X, Y, W, H int
	str        string
}

type Game struct {
	Buttons     []Button
	user_choice string
	comp_choice string
	result      string
}

func main() {

	ebiten.SetWindowTitle("Rock, Paper, Scissors - Second")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	buttons := []Button{}

	labels := []string{"Rock", "Paper", "Scissors"}

	for i := 0; i < 3; i++ {
		buttons = append(buttons, Button{
			X:   (10 + ((buttonWidth + 10) * i)),
			Y:   40,
			W:   buttonWidth,
			H:   buttonHeight,
			str: labels[i],
		})
	}

	game := &Game{
		Buttons:     buttons,
		user_choice: "",
		comp_choice: "",
		result:      "",
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {

	g.ButtonPressed()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// DRAWING BUTTONS
	for _, button := range g.Buttons {

		vector.DrawFilledRect(screen, float32(button.X), float32(button.Y), float32(buttonWidth), float32(buttonHeight), color.White, false)

		text.Draw(screen, button.str, game_font, button.X+20, button.Y+25, color.Black)
	}

	user_str := fmt.Sprintf("User's Choice: %v", g.user_choice)
	comp_str := fmt.Sprintf("Comp's Choice: %v", g.comp_choice)
	res := fmt.Sprintf("Result: %v", g.result)

	text.Draw(screen, user_str, game_font, 10, 125, color.White)
	text.Draw(screen, comp_str, game_font, 10, 175, color.White)
	text.Draw(screen, res, game_font, 10, 225, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) GenerateCompChoice() {

	choices := []string{"Rock", "Paper", "Scissors", "Rock", "Scissors", "Paper", "Paper", "Rock", "Scissors"}

	g.comp_choice = choices[generateRandom()]
}

func (g *Game) ButtonPressed() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		x_pos, y_pos := ebiten.CursorPosition()

		for _, button := range g.Buttons {

			if x_pos >= button.X && x_pos <= button.X+buttonWidth {
				if y_pos >= button.Y && y_pos <= button.Y+buttonHeight {

					g.user_choice = button.str
					g.GenerateCompChoice()

					dicts := map[string]string{
						"Rock":     "Paper",
						"Paper":    "Scissors",
						"Scissors": "Rock",
					}

					comp_won := false

					if g.user_choice == g.comp_choice {
						g.result = "It's a Draw!"
					} else {

						for k, v := range dicts {

							if k == g.user_choice && v == g.comp_choice {
								g.result = "Comp Wins!"
								comp_won = true
							}
						}

						if !comp_won {
							g.result = "Player Wins!"
						}
					}

				}
			}
		}
	}
}
