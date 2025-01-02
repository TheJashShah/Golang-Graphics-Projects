package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
)

type node struct {
	position int
	next     *node
}

type Cell struct {
	X, Y, W, H int
	marker     string
	filled     bool
}

type Reset struct {
	Cell
}

func CreateNode(pos int) *node {

	return &node{
		position: pos,
		next:     nil,
	}
}

func enqueue(pos int, head *node) *node {

	if head.position == -1 {
		head.position = pos
	} else {

		new := CreateNode(pos)
		temp := head

		for temp.next != nil {
			temp = temp.next
		}

		temp.next = new
	}

	return head
}

func dequeue(head *node) (int, *node) {

	if head.position == -1 {
		return -1, head
	}

	temp := head
	head = head.next

	return temp.position, head
}

func getLength(head *node) int {

	if head.position == -1 {
		return 0
	}

	len := 0

	temp := head

	for temp != nil {
		len++
		temp = temp.next
	}

	return len
}

func (g *Game) isCellClicked() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		X, Y := ebiten.CursorPosition()

		for i := range g.cells {

			cell := &g.cells[i]

			if (X >= cell.X && X <= cell.X+cell.W) && (Y >= cell.Y && Y <= cell.Y+cell.H) {

				if !cell.filled {

					if g.clicks%2 == 0 {
						cell.marker = "X"
						g.player_1_queue = enqueue(i, g.player_1_queue)
						g.board[i+1] = "X"
					} else {
						cell.marker = "O"
						g.player_2_queue = enqueue(i, g.player_2_queue)
						g.board[i+1] = "O"
					}
					g.clicks++
					cell.filled = true

				}
			}
		}
	}
}

func (g *Game) handleMarkerRemoval() {

	var pos int

	if getLength(g.player_1_queue) == 4 {

		pos, g.player_1_queue = dequeue(g.player_1_queue)
		g.cells[pos].marker = ""
		g.board[pos+1] = ""
		g.cells[pos].filled = false

	}

	if getLength(g.player_2_queue) == 4 {

		pos, g.player_2_queue = dequeue(g.player_2_queue)
		g.cells[pos].marker = ""
		g.board[pos+1] = ""
		g.cells[pos].filled = false
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

func (g *Game) ResetClicked() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		X, Y := ebiten.CursorPosition()

		if (X >= g.reset.X && X <= g.reset.X+g.reset.W) && (Y >= g.reset.Y && Y <= g.reset.Y+g.reset.H) {

			g.reset_visible = false

			for i := range g.board {
				g.board[i] = ""
			}

			for i := range g.cells {
				g.cells[i].marker = ""
				g.cells[i].filled = false
			}

			g.player_1_queue = &node{
				position: -1,
				next:     nil,
			}

			g.player_2_queue = &node{
				position: -1,
				next:     nil,
			}

			g.clicks = 0
		}
	}
}

func (g *Game) CheckWin(marker string) bool {

	win := false

	if (g.board[1] == marker && g.board[1] == g.board[2] && g.board[2] == g.board[3]) ||
		(g.board[4] == marker && g.board[4] == g.board[5] && g.board[5] == g.board[6]) ||
		(g.board[7] == marker && g.board[7] == g.board[8] && g.board[8] == g.board[9]) ||
		(g.board[1] == marker && g.board[1] == g.board[4] && g.board[4] == g.board[7]) ||
		(g.board[2] == marker && g.board[2] == g.board[5] && g.board[5] == g.board[8]) ||
		(g.board[3] == marker && g.board[3] == g.board[6] && g.board[6] == g.board[9]) ||
		(g.board[1] == marker && g.board[1] == g.board[5] && g.board[5] == g.board[9]) ||
		(g.board[3] == marker && g.board[3] == g.board[5] && g.board[5] == g.board[7]) {
		win = true
	}

	return win
}

func (g *Game) UpdateStrings() {

	if !g.reset_visible {

		if g.clicks%2 == 0 {
			g.display_str = "It's X's Turn!"
		} else {
			g.display_str = "It's O's Turn!"
		}
	} else {

		if g.CheckWin("X") {
			g.display_str = "Player with X Wins!"
		} else {
			g.display_str = "Player with O Wins!"
		}
	}
}
