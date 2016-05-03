package main

import (
	"fmt"
	"math"
//	"math/rand"
	"runtime"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
	"log"

	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	empty = 0
	player_one = 1
	player_two = -player_one
	searchMaxTime = 500000000 * time.Nanosecond
	searchMaxDepth = 20 
)

const (
	VerticalAxis = 1 << iota
	HorizontalAxis
	LeftDiagAxis	// haut droite
	RightDiagAxis	// bas droite
)

type mustdo struct {
	Todo	bool
	X, Y    int
}

type Board [19][19]int
type FreeThreesAxis [2][19][19]int

// copy[0] "score" ia
// copy[1] "score" player 
// copy[2] "capturable"
// copy[3] forbiden ia
// copy[4] forbiden player
// TODO: Make it an object
type BoardData [19][19][5]int

const (
	winTitle string = "Go-Gomoku"
	winTitleDebug string = "Go-Debug"
	winWidth, winHeight int = 800, 880
)

var victory mustdo

func isInBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func isValidMove(board *Board, freeThrees *[2]Board, x, y, player int) bool {
	return isInBounds(x, y) &&
		board[y][x] == empty &&
		!doesDoubleFreeThree(freeThrees, x, y, player)
}

func checkVictory(values *Board, nb int, y int, x int) bool {
	f := func (incx, incy int) int {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !isInBounds(x, y) || values[y][x] != nb {
				return i
			}
			x += incx
			y += incy
		}
		return 5
	}
	if 	(f(-1, -1) + f(1, 1) >= 4 && !checkCaptures(values, nb, x, y, 1, 1)) ||
		(f(1, -1) + f(-1, 1) >= 4 && !checkCaptures(values, nb, x, y, 1, -1)) ||
		(f(0, -1) + f(0, 1) >= 4 && !checkCaptures(values, nb, x, y, 0, 1)) ||
		(f(-1, 0) + f(1, 0) >= 4 && !checkCaptures(values, nb, x, y, 1, 0)) {
		return true
	}
	return false
}

// TODO: Create a list of possible anti-victory captures
func checkCaptures(values *Board, nb, x, y, incx, incy int) bool {
	checkAxis := func (x, y, incx, incy int) bool {
		if !isInBounds(x - incx, y - incy) || !isInBounds(x + 2 * incx, y + 2 * incy) {
			return false
		}
		if values[y + incy][x + incx] == nb {
			if values[y + 2 * incy][x + 2 * incx] == -nb && values[y - incy][x - incx] == 0 {
				victory.X = x - incx
				victory.Y = y - incy
				victory.Todo = true
				return true
			} else if values[y + 2 * incy][x + 2 * incx] == 0   && values[y - incy][x - incx] == -nb {
				victory.X = x + 2 * incx
				victory.Y = y + 2 * incy
				victory.Todo = true
				return true
			}
		}
		return false
	}
	f := func (incx, incy int) bool {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !isInBounds(x, y) || values[y][x] != nb {
				return false
			} else if checkAxis(x, y, -1, -1) || checkAxis(x, y, 1, 1) ||
				checkAxis(x, y, 1, -1) || checkAxis(x, y, -1, 1) ||
				checkAxis(x, y, 0, -1) || checkAxis(x, y, 0, 1) ||
				checkAxis(x, y, -1, 0) || checkAxis(x, y, 1, 0) {
				return true
			}
			x += incx
			y += incy
		}
		return false
	}
	return f(incx, incy) || f(-incx, -incy)
}

func checkDoubleThree(board, freeThrees *Board, x, y, color int) {
	// TODO: SOOOOO SLOOOOOW do something

	defer timeFunc(time.Now(), "checkDoubleThree")

	const (
		pat1 = 0x1A5 // -00--
		pat2 = 0x199 // -0-0-
		pat3 = 0x169 // --00-
		mask = 0x3FF // 2 * 5 bits
	)

	checkAxis := func(x, y, incx, incy, axis int) {
		if !isInBounds(x, y) || board[y][x] != empty {
			return
		}
		flags := uint32(0)

		// TODO: Perform this formatting on the entire line/column/diagonal
		tmp_x, tmp_y := x - incx*4, y - incy*4
		i := uint(0)
		for ; i < 4; i++ {
			if isInBounds(tmp_x, tmp_y) {
				flags |= uint32(board[tmp_y][tmp_x] * color + 1) << ((7 - i)*2)
			}
			tmp_x, tmp_y = tmp_x + incx, tmp_y + incy
		}
		tmp_x, tmp_y = x + incx, y + incy
		for ; i < 8; i++ {
			if isInBounds(tmp_x, tmp_y) {
				flags |= uint32(board[tmp_y][tmp_x] * color + 1) << ((7 - i)*2)
			} else {
				break
			}
			tmp_x, tmp_y = tmp_x + incx, tmp_y + incy
		}

		pos1 := (flags >> (2*3)) & mask
		pos2 := (flags >> (2*2)) & mask
		pos3 := (flags >> (2*1)) & mask
		pos4 := flags & mask
		createsFreeThree := pos4 == pat1 ||
							pos4 == pat2 ||
							pos4 == pat3 ||
							pos3 == pat1 ||
							pos3 == pat2 ||
							pos3 == pat3 ||
							pos2 == pat1 ||
							pos2 == pat2 ||
							pos2 == pat3 ||
							pos1 == pat1 ||
							pos1 == pat2 ||
							pos1 == pat3
		if createsFreeThree {
			freeThrees[y][x] |= axis
		} else {
			freeThrees[y][x] &= ^axis
		}
	}

	if board[y][x] == empty {
		checkAxis(x, y, 0, 1, VerticalAxis)
		checkAxis(x, y, 1, 0, HorizontalAxis)
		checkAxis(x, y, 1, -1, LeftDiagAxis)
		checkAxis(x, y, 1, 1, RightDiagAxis)
	}
	for i := 1; i <= 4; i++ {
		checkAxis(x, y + i, 0, 1, VerticalAxis)
		checkAxis(x, y - i, 0, 1, VerticalAxis)
		checkAxis(x + i, y, 1, 0, HorizontalAxis)
		checkAxis(x - i, y, 1, 0, HorizontalAxis)
		checkAxis(x + i, y - i, 1, -1, LeftDiagAxis)
		checkAxis(x - i, y + i, 1, -1, LeftDiagAxis)
		checkAxis(x + i, y + i, 1, 1, RightDiagAxis)
		checkAxis(x - i, y - i, 1, 1, RightDiagAxis)
	}
}

func getCaptures(board *Board, y, x, player  int, captures *[]Position) {
	captureOnAxis := func (incx, incy int) {
		if !isInBounds(x + 3 * incx, y + 3 * incy) {
			return
		}
		if	board[y + incy][x + incx] == -player &&
			board[y + 2 * incy][x + 2 * incx] == -player &&
		 	board[y + 3 * incy][x + 3 * incx] == player {
			*captures = append(*captures, Position{x + incx, y + incy})
			*captures = append(*captures, Position{x + incx * 2, y + incy * 2})
		}
	}
	captureOnAxis(-1, -1)
	captureOnAxis(1, 1)
	captureOnAxis(1, -1)
	captureOnAxis(-1, 1)
	captureOnAxis(0, -1)
	captureOnAxis(0, 1)
	captureOnAxis(-1, 0)
	captureOnAxis(1, 0)
}

func doCaptures(board *Board, captures *[]Position) {
	for _, capture := range *captures {
		board[capture.y][capture.x] = empty
	}
}

func undoCaptures(board *Board, captures *[]Position, player int) {
	for _, capture := range *captures {
		board[capture.y][capture.x] = -player
	}
}

func doesDoubleFreeThreePlayer(freeThrees *Board, x, y int) bool {
	freeThreesCount := 0
	point := freeThrees[y][x]
	if (point == 0) {
		return false
	}
	for i := uint(0); i < 4; i++ {
		if (point & (1 << i)) != 0 {
			freeThreesCount++
		}
	}
	return freeThreesCount == 2
}

func doesDoubleFreeThree(freeThrees *[2]Board, x, y, player int) bool {
	playerId := (player + 1) / 2
	return doesDoubleFreeThreePlayer(&freeThrees[playerId], x, y)
}

func updateFreeThrees(board *Board, freeThrees *[2]Board, x, y, player int, captures []Position) {
	// TODO: Two functions -> update from move and update from move cancelation
	defer timeFunc(time.Now(), "updateFreeThree")

	if (board[y][x] != empty) {
		freeThrees[0][y][x] = 0
		freeThrees[1][y][x] = 0
	}
	checkDoubleThree(board, &freeThrees[(player + 1) / 2], x, y, player)
	checkDoubleThree(board, &freeThrees[(-player + 1) / 2], x, y, -player)
	for _, pos := range captures {
		if (board[pos.y][pos.x] != empty) {
			freeThrees[0][pos.y][pos.x] = 0
			freeThrees[1][pos.y][pos.x] = 0
		}
		checkDoubleThree(board, &freeThrees[(player + 1) / 2], pos.x, pos.y, player)
		checkDoubleThree(board, &freeThrees[(-player + 1) / 2], pos.x, pos.y, -player)
	}
}

func checkRules(values *Board, freeThrees *[2]Board, capture *[3]int, x, y, player int) int {
	if doesDoubleFreeThree(freeThrees, x, y, player) {
		fmt.Printf("Nope\n")
		return player
	}
	values[y][x] = player
	victory := checkVictory(values, player, y, x)
	if victory == true {
		fmt.Printf("Victorye \\o/ %d\n", player)
		return 0
	}
	captures := make([]Position, 0, 16)
	getCaptures(values, y, x, player, &captures)
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
	} else if t > 18{
		t = 18
	}
	return t
}

func init() {
    runtime.LockOSThread()
}

func evaluateAllBoard(player int, value *Board, better *BoardData, capture *[3]int) {
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			evaluateBoard(value, x, y, player, better, capture)
		}
	}
}

func run() int {
	var event sdl.Event
	var running bool
	var err error
	var player int
	victory.Todo = false
	var capture [3]int
	var values Board
	var freeThrees [2]Board
	var better BoardData

	f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %s\n", err)
		return 1
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("---NEW GAME---\n")

	sdl.Init(sdl.INIT_EVERYTHING)
	if err := ttf.Init(); err != nil {
		fmt.Println(err)
		return 3
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow(winTitle, 800, 0,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	windowb, err := sdl.CreateWindow(winTitleDebug, 0, 0,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer windowb.Destroy()
	rendererb, err := sdl.CreateRenderer(windowb, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer rendererb.Destroy()

	_= rendererb.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	running = true
	player = 1
	var px, py int
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				//fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)

				/*
				if player == player_one && t.Type == 1025 {
				/*/
					if t.Type == 1025 {
				//*/

					py = mousePositionToGrid(float64(t.Y))
					px = mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", px, py)
					log.Printf("p1 -> X |%3d| Y|%3d|\n", px, py)
					if victory.Todo == true {
						if px == victory.X && py == victory.Y {
							player = checkRules(&values, &freeThrees, &capture, px, py, player)
							victory.Todo = false
						} else {
							fmt.Printf("you must play in [%d][%d]\n", victory.X, victory.Y)
						}
						fmt.Println(values)
					} else if values[py][px] == 0 {
						player = checkRules(&values, &freeThrees, &capture, px, py, player)
					}
					evaluateAllBoard(player, &values, &better, &capture)
				}
			case *sdl.KeyUpEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}

		if player == player_two {
			if victory.Todo == true {
				fmt.Printf("IA must play -> x[%d] y [%d]\n", victory.X, victory.Y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", victory.X, victory.Y)
				player = checkRules(&values, &freeThrees, &capture, victory.X, victory.Y, player)
				victory.Todo = false
			} else {
				var x, y int
				x, y, better = search(&values, &freeThrees, player, px, py, 4, &capture)
				fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", x, y)
				if values[y][x] == 0 {
					player = checkRules(&values, &freeThrees, &capture, x, y, player)
				}
			}
			displayAverages()
			resetTimer()
		}

		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		drawGrid(renderer)
		drawClic(renderer, &values, &capture, &freeThrees)
		renderer.Present()

		_ = rendererb.SetDrawColor(236, 240, 241, 0)
		rendererb.Clear()
		drawGrid(rendererb)
		drawClic(rendererb, &values, &capture, &freeThrees)
		draweval(rendererb, &better)
		rendererb.Present()
	}
	return 0
}

func main() {
	fmt.Printf("%d %d", HorizontalAxis, LeftDiagAxis)
	os.Exit(run())
}
