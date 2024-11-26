package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	//"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 800
	screenHeight = 800
	buttonWidth  = 40
	buttonHeight = 50
	resetWidth   = 150
	resetHeight  = 30

	padding = 16
)

var (
	game_font   font.Face
	button_font font.Face
	str_font    font.Face

	img_1 *ebiten.Image
	img_2 *ebiten.Image
	img_3 *ebiten.Image
	img_4 *ebiten.Image
	img_5 *ebiten.Image
	img_6 *ebiten.Image
	img_7 *ebiten.Image

	words = [...]string{"CRISTIANO", "MESSI", "RONALDO", "VINICIUS", "NEYMAR", "SALAH", "MBAPPE", "DEBRUYNE", "RAMOS", "BENZEMA", "GARETH", "MARCELO", "ZIDANE", "ANCELOTTI", "YAMAL", "MODRIC"}
)

type Button struct {
	X, Y, W, H int
	letter     string
	visible    bool
}

type Reset struct {
	X, Y, W, H int
	visible    bool
}

type Game struct {
	buttons          []Button
	word             string
	lives            int
	guess            []string
	won_str_visible  bool
	lost_str_visible bool
	imgs             []*ebiten.Image
	reset            Reset
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

	button_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    22,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	game_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	str_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	img_1, _, err = ebitenutil.NewImageFromFile("assets/hangman0.png")

	if err != nil {
		log.Fatal(err)
	}
	img_2, _, err = ebitenutil.NewImageFromFile("assets/hangman1.png")

	if err != nil {
		log.Fatal(err)
	}
	img_3, _, err = ebitenutil.NewImageFromFile("assets/hangman2.png")

	if err != nil {
		log.Fatal(err)
	}
	img_4, _, err = ebitenutil.NewImageFromFile("assets/hangman3.png")

	if err != nil {
		log.Fatal(err)
	}
	img_5, _, err = ebitenutil.NewImageFromFile("assets/hangman4.png")

	if err != nil {
		log.Fatal(err)
	}
	img_6, _, err = ebitenutil.NewImageFromFile("assets/hangman5.png")

	if err != nil {
		log.Fatal(err)
	}
	img_7, _, err = ebitenutil.NewImageFromFile("assets/hangman6.png")

	if err != nil {
		log.Fatal(err)
	}
}

func stringBuilder(len int) []string {

	str := []string{}

	for i := 0; i < len; i++ {
		str = append(str, "_ ")
	}

	return str
}

func generateRandomNum(limit int) int {

	return rand.Intn(limit)
}

func main() {

	ebiten.SetWindowTitle("Hangman - Third")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	buttons := []Button{}

	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	str_list := strings.Split(str, "")

	for i := 0; i < 14; i++ {

		buttons = append(buttons, Button{
			X:       padding + ((buttonWidth + padding) * i),
			Y:       680,
			W:       buttonWidth,
			H:       buttonHeight,
			letter:  str_list[i],
			visible: true,
		})
	}

	for j := 0; j < 12; j++ {

		buttons = append(buttons, Button{
			X:       56 + ((buttonWidth + padding) * j),
			Y:       740,
			W:       buttonWidth,
			H:       buttonHeight,
			letter:  str_list[j+14],
			visible: true,
		})
	}

	reset := Reset{
		X:       (screenWidth - resetWidth) / 2,
		Y:       650,
		W:       resetWidth,
		H:       resetHeight,
		visible: false,
	}

	word := words[generateRandomNum(len(words))]

	guess := stringBuilder(len(word))
	lives := 7

	game := &Game{
		buttons:          buttons,
		word:             word,
		lives:            lives,
		guess:            guess,
		won_str_visible:  false,
		lost_str_visible: false,
		imgs:             []*ebiten.Image{img_1, img_2, img_3, img_4, img_5, img_6, img_7},
		reset:            reset,
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}

}

func returnWidth(fn font.Face, message string) int {

	width := 0

	for _, r := range message {

		r_w, _ := fn.GlyphAdvance(r)

		width += r_w.Ceil()
	}

	return width
}

func (g *Game) Update() error {

	g.isButtonClicked()
	g.determineEnd()
	g.isResetClicked()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.White)

	for _, button := range g.buttons {

		if button.visible {
			vector.DrawFilledRect(screen, float32(button.X), float32(button.Y), float32(buttonWidth), float32(buttonHeight), color.Black, false)
			vector.DrawFilledRect(screen, float32(button.X+2), float32(button.Y+2), float32(buttonWidth-4), float32(buttonHeight-4), color.White, false)
			text.Draw(screen, button.letter, button_font, button.X+10, button.Y+30, color.Black)
		}
	}

	guess_word := strings.Join(g.guess, " ")
	guessStr := fmt.Sprintf("Your Word is: %v", guess_word)
	text.Draw(screen, guessStr, game_font, 10, 100, color.Black)

	livesStr := fmt.Sprintf("You have %d lives remaining", g.lives)
	text.Draw(screen, livesStr, game_font, 10, 160, color.Black)

	if g.won_str_visible {
		won_str := "You have correctly Guessed the Word!"

		width := returnWidth(str_font, won_str)
		height := str_font.Metrics().Height.Ceil()

		x, y := (screenWidth-width)/2, (screenHeight-height)/2+height
		text.Draw(screen, won_str, str_font, x, y, color.Black)
	}

	if g.lost_str_visible {
		lost_str := fmt.Sprintf("The Correct Word was: %v", g.word)

		width := returnWidth(button_font, lost_str)
		height := str_font.Metrics().Height.Ceil()

		x, y := (screenWidth-width)/2, (screenHeight-height)/2+height
		text.Draw(screen, lost_str, button_font, x, y, color.Black)
	}

	if g.lives <= 6 && g.lives >= 0 {

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(590, 300)

		screen.DrawImage(g.imgs[len(g.imgs)-g.lives-1], options)
	}

	if g.reset.visible {

		vector.DrawFilledRect(screen, float32(g.reset.X), float32(g.reset.Y), float32(g.reset.W), float32(g.reset.H), color.Black, false)
		vector.DrawFilledRect(screen, float32(g.reset.X+2), float32(g.reset.Y+2), float32(g.reset.W-4), float32(g.reset.H-4), color.White, false)

		str := "RESET"

		width := returnWidth(button_font, str)
		height := button_font.Metrics().Height.Ceil()

		x, y := (screenWidth-width)/2, 650+height

		text.Draw(screen, str, button_font, x, y, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return screenWidth, screenHeight
}

func (g *Game) isButtonClicked() {

	indices := []int{}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		x_pos, y_pos := ebiten.CursorPosition()

		for i := range g.buttons {

			button := &g.buttons[i]

			if button.visible {

				if (x_pos >= button.X && x_pos <= button.X+buttonWidth) &&
					(y_pos >= button.Y && y_pos <= button.Y+buttonHeight) {
					button.visible = false

					if strings.Contains(g.word, button.letter) {

						for i := range g.word {
							if (string(g.word[i])) == button.letter {
								indices = append(indices, i)
							}
						}

						for _, i := range indices {
							g.guess[i] = button.letter
						}
					} else {
						g.lives--
					}
				}
			}
		}
	}
}

func (g *Game) determineEnd() {

	if g.lives >= 1 && strings.Join(g.guess, "") == g.word {

		for i := range g.buttons {
			button := &g.buttons[i]
			button.visible = false
		}

		g.won_str_visible = true
		g.reset.visible = true

	} else if g.lives <= 0 {

		for i := range g.buttons {
			button := &g.buttons[i]
			button.visible = false
		}

		g.lost_str_visible = true
		g.reset.visible = true
	}
}

func (g *Game) isResetClicked() {

	if g.reset.visible {

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

			x, y := ebiten.CursorPosition()

			if (x >= g.reset.X && x <= g.reset.X+g.reset.W) && (y >= g.reset.Y && y <= g.reset.Y+g.reset.H) {

				for i := range g.buttons {
					button := &g.buttons[i]
					button.visible = true
				}

				g.lost_str_visible = false
				g.won_str_visible = false

				g.lives = 7

				g.word = words[generateRandomNum(len(words))]
				g.guess = stringBuilder(len(g.word))

				g.reset.visible = false
			}
		}
	}
}
