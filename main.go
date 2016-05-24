// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: tmielcza <marvin@42.fr>                    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2016/05/16 18:08:05 by tmielcza          #+#    #+#             //
//   Updated: 2016/05/23 21:00:09 by tmielcza         ###   ########.fr       //
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

type AlignScore struct {
	score_player_one	int
	score_player_two	int
	x, y				int
}

type Board [19][19]int
type FreeThreesAxis [2][19][19]int

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

var textDrawer *TextDrawer

func isInBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func isValidMove(board *Board, freeThrees *[2]Board, x, y, player int) bool {
	return isInBounds(x, y) &&
		board[y][x] == empty &&
		!doesDoubleFreeThree(freeThrees, x, y, player)
}

type MoveType int

const (
	regularMove MoveType = iota
	winByAlignment
	winByCapture
)

func canPlay(board *Board, freeThrees *[2]Board, forcedCaptures []Position, x, y, player int) bool {
	return board[y][x] == empty &&
		!doesDoubleFreeThree(freeThrees, x, y, player) &&
		(forcedCaptures == nil || containsPosition(forcedCaptures, Position{x, y}))
}

func checkRules(board *Board, freeThrees, alignTable *[2]Board, capturesNb *[3]int, x, y, player int) (MoveType, []Position) {
//	a,b,c,d := getBestScore(board, alignTable, x, y, player)
//	fmt.Printf("best[%d,%d] worst[%d,%d]\n", a,b,c,d)
	updateAlign(board, alignTable, x, y, player)
	board[y][x] = player
	alignmentType, forcedCaptures := checkVictory(board, capturesNb, x, y, player)
	switch alignmentType {
	case winningAlignment:
		return winByAlignment, nil
	case regularAlignment:
		forcedCaptures = nil
	}
	captures := make([]Position, 0, 16)
	getCaptures(board, x, y, player, &captures)
	doCaptures(board, &captures)
	capturesNb[player + 1] += len(captures)
	updateFreeThrees(board, freeThrees, x, y, player, captures)
	if capturesNb[player + 1] >= 10 {
		return winByCapture, nil
	}
	return regularMove, forcedCaptures
}

func containsPosition(captures []Position, pos Position) bool {
	for _, capt := range captures {
		if pos == capt {
			return true
		}
	}
	return false
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
		searchTime		time.Duration
	)

	var (
		player, px, py	int
		capture			[3]int
		values			Board
		freeThrees		[2]Board
		alignTable		[2]Board
		better			BoardData
		forcedCaptures	[]Position = nil
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

	// loop
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				//fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if  t.Type == 1025 && ((player_mode == 1 && player == player_one) || (player_mode == 2)) {
					py = mousePositionToGrid(float64(t.Y))
					px = mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", px, py)
					log.Printf("p1 -> X |%3d| Y|%3d|\n", px, py)
					if canPlay(&values, &freeThrees, forcedCaptures, px, py, player) {

						moveType, newForcedCaptures := checkRules(&values, &freeThrees, &alignTable, &capture, px, py, player)
						forcedCaptures = newForcedCaptures
						if moveType != regularMove {
							player = 0
						}
						player = -player
					} else {
						fmt.Println("Can't play here.")
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
					log.Printf("---NEW GAME---\n")
					player = 1
				}
				// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%d\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}

		// IA
		if player_mode == 1 && player == player_two {

			var x, y int
			startTime := time.Now()
			x, y, better = search(&values, &freeThrees, &alignTable, player, px, py, 5, &capture, forcedCaptures)
			searchTime = time.Since(startTime)
			fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
			log.Printf("IA -> X |%3d| Y|%3d|\n", x, y)
			if canPlay(&values, &freeThrees, forcedCaptures, x, y, player) {
				moveType, newForcedCaptures := checkRules(&values, &freeThrees, &alignTable, &capture, x, y, player)
				forcedCaptures = newForcedCaptures
				if moveType != regularMove {
					player = 0
				}
				player = -player
			} else {
				fmt.Println("Can't play here.")
			}
			displayAverages()
			resetTimer()
		}

		// Rendering Window(s)
		if player_mode > 0 {
			_ = renderer.SetDrawColor(236, 240, 241, 0)
			renderer.Clear()
			drawClic(renderer, &values, &capture, forcedCaptures, &freeThrees, int(searchTime))
			if player == 0 {
				drawRestartPanel(renderer)
			}
			renderer.Present()
			if debug == true {
				_ = rendererb.SetDrawColor(236, 240, 241, 0)
				rendererb.Clear()
				drawClic(rendererb, &values, &capture, forcedCaptures, &freeThrees, int(searchTime))
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
