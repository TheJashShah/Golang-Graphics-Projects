package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// CONSTANTS
const (
	screenWidth  = 640
	screenHeight = 480
	paddleSpeed  = 6
	ballSpeed    = 4

	paddleWidth  = 20
	paddleHeight = 100
)

// GENERIC OBJECT STRUCT
type Object struct {
	X, Y, W, H int
}

// PADDLE
type Paddle struct {
	Object
}

// BALL
type Ball struct {
	Object
	x_speed int
	y_speed int
}

// GAME
type Game struct {
	paddle_1  Paddle
	paddle_2  Paddle
	ball      Ball
	score     int
	highScore int
}

func main() {

	ebiten.SetWindowTitle("Pong - The First")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paddle_1 := Paddle{
		Object: Object{
			X: 20,
			Y: 190,
			W: paddleWidth,
			H: paddleHeight,
		},
	}

	paddle_2 := Paddle{
		Object: Object{
			X: 600,
			Y: 190,
			W: paddleWidth,
			H: paddleHeight,
		},
	}

	ball := Ball{
		Object: Object{
			X: 310,
			Y: 230,
			W: 20,
			H: 20,
		},
		x_speed: ballSpeed,
		y_speed: ballSpeed,
	}

	game := &Game{
		paddle_1:  paddle_1,
		paddle_2:  paddle_2,
		ball:      ball,
		score:     0,
		highScore: 0,
	}

	err := ebiten.RunGame(game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {

	// PADDLE MOVEMENTS
	g.paddle_1.MoveonKeyPress()
	g.paddle_2.MoveonKeyPress()

	// BALL FUNCTIONS
	g.ball.move()
	g.ball.collisionWithWalls(g)

	//COLLISION WITH PADDLES
	g.collisionWithPaddles()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// PADDLE_1
	vector.DrawFilledRect(screen, float32(g.paddle_1.X), float32(g.paddle_1.Y), float32(g.paddle_1.W), float32(g.paddle_1.H), color.White, false)

	// PADDLE_2
	vector.DrawFilledRect(screen, float32(g.paddle_2.X), float32(g.paddle_2.Y), float32(g.paddle_2.W), float32(g.paddle_2.H), color.White, false)

	// BALL
	vector.DrawFilledRect(screen, float32(g.ball.X), float32(g.ball.Y), float32(g.ball.W), float32(g.ball.H), color.White, false)

	// SCORE AND HIGHSCORE TEXTS
	scoreStr := fmt.Sprintf("Score: %v", g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 10, color.RGBA{100, 200, 250, 1})

	highscoreStr := fmt.Sprintf("High Score: %v", g.highScore)
	text.Draw(screen, highscoreStr, basicfont.Face7x13, 10, 30, color.RGBA{150, 200, 200, 1})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return screenWidth, screenHeight
}

func (p *Paddle) MoveonKeyPress() {

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.Y >= 0 {
		p.Y -= paddleSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && (p.Y <= screenHeight-p.H) {
		p.Y += paddleSpeed
	}
}

func (b *Ball) move() {

	b.X += b.x_speed
	b.Y += b.y_speed
}

func (g *Game) reset() {

	g.ball.X = 310
	g.ball.Y = 230

	g.score = 0

}

func (b *Ball) collisionWithWalls(g *Game) {

	if b.X <= ballSpeed {
		g.reset()
	}

	if b.X >= screenWidth-b.W-ballSpeed {
		g.reset()
	}

	if b.Y <= ballSpeed {
		b.y_speed = -b.y_speed
	}

	if b.Y >= screenHeight-b.H-ballSpeed {
		b.y_speed = -b.y_speed
	}
}

func (g *Game) collisionWithPaddles() {

	//PADDLE 1
	if g.ball.X <= g.paddle_1.X+paddleWidth &&
		(g.ball.X+g.ball.W >= g.paddle_1.X) &&
		(g.ball.Y+g.ball.H >= g.paddle_1.Y) &&
		(g.ball.Y <= (g.paddle_1.Y + paddleHeight)) {
		g.ball.X = g.paddle_1.X + paddleWidth + 1
		g.ball.x_speed = -g.ball.x_speed
		g.score++
	}

	//PADDLE 2
	if (g.ball.X+g.ball.W >= g.paddle_2.X) &&
		(g.ball.X <= g.paddle_2.X+paddleWidth) &&
		(g.ball.Y <= g.paddle_2.Y+paddleHeight) &&
		(g.ball.Y+g.ball.H >= g.paddle_2.Y) {
		g.ball.X = g.paddle_2.X - paddleWidth - 1
		g.ball.x_speed = -g.ball.x_speed

		g.score++
	}

	//HighScore Calculation
	if g.score > g.highScore {
		g.highScore = g.score
	}
}
