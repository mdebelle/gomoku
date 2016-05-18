// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: tmielcza <marvin@42.fr>                    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2016/05/16 18:08:05 by tmielcza          #+#    #+#             //
//   Updated: 2016/05/18 19:28:30 by tmielcza         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"math"
	"runtime"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
	"log"

	"github.com/veandco/go-sdl2/sdl_ttf"
)

const debug = false

const (
	empty = 0
	player_one = 1
	player_two = -player_one
	searchMaxTime = 500000000 * time.Nanosecond
	searchMaxDepth = 20
)

const (
	VerticalAxisMask = 1 << iota
	HorizontalAxisMask
	LeftDiagAxisMask	// haut droite
	RightDiagAxisMask	// bas droite
)

const (
	VerticalAxis = iota
	HorizontalAxis
	LeftDiagAxis	// haut droite
	RightDiagAxis	// bas droite
)

type Position struct {
	x, y int
}

type mustdo struct {
	Todo	bool
	X, Y    int
}

type AlignScore struct {
	score_player_one	int
	score_player_two	int
	x, y				int
}

type Board [19][19]int
type FreeThreesAxis [2][19][19]int

// copy[0] "score" ia
// copy[1] "score" player 
// copy[2] "capturable"
// copy[3] forbiden ia
// copy[4] forbiden player
// korpy[5] score
// porky[6] isTested?
// plwoer[7] base Score
// TODO: Make it an object (or not..)
type BoardData [19][19][8]int

const (
	winTitle string = "Go-Gomoku"
	winTitleDebug string = "Go-Debug"
	winWidth, winHeight int = 800, 880
)

type AlignmentType int

const (
	regularAlignment AlignmentType = iota
	capturableAlignment
	winningAlignment
)

type Unit struct {}

var victory mustdo

var textDrawer *TextDrawer

func isInBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func isValidMove(board *Board, freeThrees *[2]Board, x, y, player int) bool {
	return isInBounds(x, y) &&
		board[y][x] == empty &&
		!doesDoubleFreeThree(freeThrees, x, y, player)
}

func updateAlignAfterCapture(board *Board, alignTable *[2]Board, lst []Position, hey int) {

	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	resetAxeScore := func(x, y, incx, incy, axe, c int) {
		var state bool

		if !isInBounds(x,y) { return }

		player := board[y][x]

		if player != 0 {
			for i := 1; i < 5; i++ {
				if isInBounds(x+(i*incx), y+(i*incy)) {
					if !state && board[y+(i*incy)][x+(i*incx)] == 0 {
						alignTable[(player+1)/2][y+(i*incy)][x+(i*incx)] |= (1 << uint(axe-i))
						clearOponentSituation(player, y+(i*incy), x+(i*incx), axe, i)
					} else if board[y+(i*incy)][x+(i*incx)] == -player {
						state = true
					}
				}
			}
		} else {
			if axe % 8 == 0 { 
				alignTable[(hey+1)/2][y][x] &= ^(1 << uint(axe-4-c))
			} else {
				alignTable[(hey+1)/2][y][x] &= ^(1 << uint(axe+4-c))
			}
		}
	}

	const (
		axeLeft = 4
		axeRight = 8
		axeTop = 12
		axeBottom = 16
		axeLeftTop = 20
		axeRightBottom = 24
		axeRightTop = 28
		axeLeftBottom = 32
	)

	for _, p := range lst {

		alignTable[(player_one + 1)/2][p.y][p.x] = 0
		alignTable[(player_two + 1)/2][p.y][p.x] = 0

		for i := 1; i < 5; i++ {
			// LeftRight
			resetAxeScore(p.x-i, p.y, 1, 0, axeRight, i)
			resetAxeScore(p.x+i, p.y, -1, 0, axeLeft, i)
			resetAxeScore(p.x, p.y-i, 0, 1, axeBottom, i)
			resetAxeScore(p.x, p.y+i, 0, -1, axeTop, i)
			resetAxeScore(p.x-i, p.y-i, 1, 1, axeRightBottom, i)
			resetAxeScore(p.x+i, p.y+i, -1, -1, axeLeftTop, i)
			resetAxeScore(p.x+i, p.y-i, -1, 1, axeLeftBottom, i)
			resetAxeScore(p.x-i, p.y+i, 1, -1, axeRightTop, i)
		}
	}
}

func checkRules(values *Board, freeThrees, alignTable *[2]Board, capture *[3]int, x, y, player int) int {
	if doesDoubleFreeThree(freeThrees, x, y, player) {
		fmt.Printf("Nope\n")
		return player
	}
	updateAlign(values, alignTable, x, y, player)
	values[y][x] = player
	victory := checkVictory(values, player, y, x)
	if victory == true {
		fmt.Printf("Victory \\o/ %d\n", player)
		return 0
	}
	captures := make([]Position, 0, 16)
	getCaptures(values, x, y, player, &captures)
	doCaptures(values, &captures)
	capture[player + 1] += len(captures)
	updateFreeThrees(values, freeThrees, x, y, player, captures)
	if capture[player + 1] >= 10 {
		fmt.Printf("capture de ouf \\o/ %d\n", player)
		return 0
	}
	return -player
}

func mousePositionToGrid(val float64) int {
	t := int(math.Floor((val - 20) / 40))
	if t < 0 {
		t = 0
	} else if t > 18 {
		t = 18
	}
	return t
}

func init() {
    runtime.LockOSThread()
}

func run() int {
	var (
		event			sdl.Event
		running			bool
		err				error
		player_mode		int
	)

	var (
		player, px, py	int
		capture			[3]int
		values			Board
		freeThrees		[2]Board
		alignTable		[2]Board
		better			BoardData
	)
	
	// Log Module
	f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %s\n", err)
		return 1
	}
	defer f.Close()
	log.SetOutput(f)
	log.Printf("---NEW GAME---\n")

	// Init SDL
	sdl.Init(sdl.INIT_EVERYTHING)
	
	// Drawing Module
	if err := ttf.Init(); err != nil {
		fmt.Println(err)
		return 3
	}
	defer ttf.Quit()
	textDrawer = NewTextDrawer()
	defer textDrawer.Dispose()

	// Main Window
	var (
		window		*sdl.Window
		renderer	*sdl.Renderer
	)
	window, err = sdl.CreateWindow(winTitle, 800, 0,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	// Debug Window
	var (
		debug		bool
		windowb		*sdl.Window
		rendererb	*sdl.Renderer
	)

	player = 1
	running = true
	victory.Todo = false

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				//fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if  (player_mode == 1 && player == player_one) || (player_mode == 2) {
					py = mousePositionToGrid(float64(t.Y))
					px = mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", px, py)
					log.Printf("p1 -> X |%3d| Y|%3d|\n", px, py)
					if victory.Todo == true {
						if px == victory.X && py == victory.Y {
							player = checkRules(&values, &freeThrees, &alignTable, &capture, px, py, player)
							victory.Todo = false
						} else {
							fmt.Printf("you must play in [%d][%d]\n", victory.X, victory.Y)
						}
						fmt.Println(values)
					} else if values[py][px] == 0 {
						a,b,c,d := getScore(&alignTable, px, py, player)
						fmt.Printf("LR %d, TB %d, LTRB %d, RTLB %d\n", a, b, c, d)
						player = checkRules(&values, &freeThrees, &alignTable, &capture, px, py, player)
					}
				}
			case *sdl.KeyUpEvent:
				// Quit Event
				if t.Keysym.Sym == 'q' || t.Keysym.Sym == 27 {
					running = false
				}

				// Select Game Mode Event
				if player_mode == 0 && (t.Keysym.Sym == '1' || t.Keysym.Sym == 1073741913) {
					player_mode = 1
				} else if player_mode == 0 && (t.Keysym.Sym == '2' || t.Keysym.Sym == 1073741914) {
					player_mode = 2
				}

				// Toggle Debug Window Event
				if player_mode > 0 && t.Keysym.Sym == 'd' {
					if debug == false {
						windowb, err = sdl.CreateWindow(winTitleDebug, 0, 0,
							winWidth, winHeight, sdl.WINDOW_SHOWN)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
							return 1
						}
						defer windowb.Destroy()
						rendererb, err = sdl.CreateRenderer(windowb, -1, sdl.RENDERER_ACCELERATED)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
							return 2
						}
						defer rendererb.Destroy()
						rendererb.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
						debug = true
					} else {
						rendererb.Destroy()
						windowb.Destroy()
						debug = false
						log.Printf("---NEW GAME---\n")
					}
				}

				// Restart option after game over Event
				if player == 0 && t.Keysym.Sym == 'a' {
					var (
						emptyvalues		Board
						emptyfreeThrees	[2]Board
						emptybetter		BoardData
					)
					capture = [3]int {0,0,0}
					values = emptyvalues
					freeThrees = emptyfreeThrees
					better =  emptybetter
					player_mode = 0
					if debug == true {
						rendererb.Destroy()
						windowb.Destroy()
						debug = false
					}
					player = 1
				}
			//	fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%d\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}

		// IA
		if player_mode == 1 && player == player_two {
			if victory.Todo == true {
				fmt.Printf("IA must play -> x[%d] y [%d]\n", victory.X, victory.Y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", victory.X, victory.Y)
				player = checkRules(&values, &freeThrees, &alignTable, &capture, victory.X, victory.Y, player)
				victory.Todo = false
			} else {
				var x, y int
				x, y, better = search(&values, &freeThrees, &alignTable, player, px, py, 5, &capture)
				fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", x, y)
				if values[y][x] == 0 {
					player = checkRules(&values, &freeThrees, &alignTable, &capture, x, y, player)
				}
			}
			displayAverages()
			resetTimer()
		}

		// Rendering Window(s)
		if player_mode > 0 {
			_ = renderer.SetDrawColor(236, 240, 241, 0)
			renderer.Clear()
			drawClic(renderer, &values, &capture, &freeThrees)
			if player == 0 {
				drawRestartPanel(renderer)
			}
			renderer.Present()
			if debug == true {
				_ = rendererb.SetDrawColor(236, 240, 241, 0)
				rendererb.Clear()
				drawClic(rendererb, &values, &capture, &freeThrees)
				draweval(rendererb, &better)
				rendererb.Present()
			}
		} else {
			drawPanel(renderer)
			if debug == true { drawPanel(rendererb) }
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
