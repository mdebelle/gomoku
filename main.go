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
		x, y := x, y
		for i := 0; i < 5; i++ {
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

	if debug { defer timeFunc(time.Now(), "checkDoubleThree") }

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
		checkAxis(x, y, 0, 1, VerticalAxisMask)
		checkAxis(x, y, 1, 0, HorizontalAxisMask)
		checkAxis(x, y, 1, -1, LeftDiagAxisMask)
		checkAxis(x, y, 1, 1, RightDiagAxisMask)
	}
	for i := 1; i <= 4; i++ {
		checkAxis(x, y + i, 0, 1, VerticalAxisMask)
		checkAxis(x, y - i, 0, 1, VerticalAxisMask)
		checkAxis(x + i, y, 1, 0, HorizontalAxisMask)
		checkAxis(x - i, y, 1, 0, HorizontalAxisMask)
		checkAxis(x + i, y - i, 1, -1, LeftDiagAxisMask)
		checkAxis(x - i, y + i, 1, -1, LeftDiagAxisMask)
		checkAxis(x + i, y + i, 1, 1, RightDiagAxisMask)
		checkAxis(x - i, y - i, 1, 1, RightDiagAxisMask)
	}
}

func getCaptures(board *Board, x, y, player  int, captures *[]Position) {
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
	if debug { defer timeFunc(time.Now(), "updateFreeThree") }

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

func getScore(alignTable *[2]Board, x, y, player int) (int, int, int, int) {

	if !isInBounds(x,y) { return 0,0,0,0 }

	const (
		maskleft = 0x0000000f
		maskright = 0x000000f0
		masktop = 0x00000f00
		maskbottom = 0x0000f000
		masklefttop = 0x000f0000
		maskrightbottom = 0x00f00000
		maskrighttop = 0x0f000000
		maskleftbottom = 0xf0000000
	)

	v := alignTable[(player+1)/2][y][x]

	axeLeftRight := (v & maskleft) + ((v & maskright) >> 4)
	axeTopBottom := ((v & masktop) >> 8) + ((v & maskbottom) >> 12)
	axeLeftTopRightBottom := ((v & masklefttop) >> 16) + ((v & maskrightbottom) >> 20)
	axeRightTopLeftBottom := ((v & maskrighttop) >> 24) + ((v & maskleftbottom) >> 28)

	return axeLeftRight, axeTopBottom, axeLeftTopRightBottom, axeRightTopLeftBottom
}

func updateAlignAfterCapture(board *Board, alignTable *[2]Board, lst []Position) {

	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	resetAxeScore := func(player, x, y, incx, incy, axe int) {
		var state bool
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
			if isInBounds(p.x-i,p.y) && board[p.y][p.x-i] != 0 {
				resetAxeScore(board[p.y][p.x-i], p.x-i, p.y, 1, 0, axeRight)
			}
			if isInBounds(p.x+i,p.y) && board[p.y][p.x+i] != 0 {
				resetAxeScore(board[p.y][p.x+i], p.x+i, p.y, -1, 0, axeLeft)
			}
			// TopBottom
			if isInBounds(p.x,p.y-i) && board[p.y-i][p.x] != 0 {
				resetAxeScore(board[p.y][p.x-i], p.x, p.y-i, 0, 1, axeBottom)
			}
			if isInBounds(p.x,p.y+i) && board[p.y+i][p.x] != 0 {
				resetAxeScore(board[p.y][p.x+i], p.x, p.y+i, 0, -1, axeTop)
			}
			// LeftTopRightBottom
			if isInBounds(p.x-i,p.y-i) && board[p.y-i][p.x-i] != 0 {
				resetAxeScore(board[p.y-i][p.x-i], p.x-i, p.y-i, 1, 1, axeRightBottom)
			}
			if isInBounds(p.x+i,p.y+i) && board[p.y+i][p.x+i] != 0 {
				resetAxeScore(board[p.y+i][p.x+i], p.x+i, p.y+i, -1, -1, axeLeftTop)
			}
			// RightTopLeftBottom
			if isInBounds(p.x+i,p.y-i) && board[p.y-i][p.x+i] != 0 {
				resetAxeScore(board[p.y-i][p.x+i], p.x+i, p.y-i, -1, 1, axeRightTop)
			}
			if isInBounds(p.x-i,p.y+i) && board[p.y+i][p.x-i] != 0 {
				resetAxeScore(board[p.y+i][p.x-i], p.x-i, p.y+i, 1, -1, axeLeftBottom)
			}

		}

	}

}

func updateAlign(board *Board, alignTable *[2]Board, x, y, player int) {

	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	updateDistanceScore := func (player, px, py, axe, i int, state bool) bool {
		if isInBounds(px, py) {
			if !state && board[py][px] == 0 {
				alignTable[(player+1)/2][py][px] |= (1 << uint(axe-i))
				clearOponentSituation(player, py, px, axe, i)
			} else if board[py][px] == -player {
				state = true
			}
		}
		return state
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

	var l, lt, t, rt, r, rb, b, lb bool

	alignTable[(player+1)/2][y][x] = 0
	alignTable[(-player+1)/2][y][x] = 0
	
	for i:= 1; i < 5; i++ {
		// Axe left right
		l = updateDistanceScore(player, x-i, y, axeLeft, i, l)
		r = updateDistanceScore(player, x+i, y, axeRight, i, r)
		// Axe top left
		t = updateDistanceScore(player, x, y-i, axeTop, i, t)
		b = updateDistanceScore(player, x, y+i, axeBottom, i, b)
		// Axe lefttop rightbottom
		lt = updateDistanceScore(player, x-i, y-i, axeLeftTop, i, lt)
		rb = updateDistanceScore(player, x+i, y+i, axeRightBottom, i, rb)
		// Axe righttop leftbottom
		rt = updateDistanceScore(player, x+i, y-i, axeRightTop, i, rt)
		lb = updateDistanceScore(player, x-i, y+i, axeLeftBottom, i, lb)
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
	var event sdl.Event
	var running bool
	var err error
	var player int
	victory.Todo = false
	var capture [3]int
	var values Board
	var freeThrees [2]Board
	var alignTable [2]Board
	var better BoardData

	var player_mode int
	var debug		bool

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

	textDrawer = NewTextDrawer()
	defer textDrawer.Dispose()

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

	var windowb *sdl.Window
	var rendererb *sdl.Renderer


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
				if  t.Type == 1025 && ((player_mode == 1 && player == player_one) || (player_mode == 2)) {
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
						getScore(&alignTable, px, py, player)
						updateAlign(&values, &alignTable, px, py, player)
						player = checkRules(&values, &freeThrees, &capture, px, py, player)
					}
				}
			case *sdl.KeyUpEvent:
				if player_mode == 0 && t.Type == 769 && (t.Keysym.Sym == '1' || t.Keysym.Sym == 1073741913) {
					player_mode = 1
				} else if player_mode == 0 && t.Type == 769 && (t.Keysym.Sym == '2' || t.Keysym.Sym == 1073741914) {
					player_mode = 2
				}
				if t.Type == 769 && (t.Keysym.Sym == 'q' || t.Keysym.Sym == 27) {
					running = false
				}
				if player_mode > 0 && t.Type == 769 && t.Keysym.Sym == 'd' {
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

						_= rendererb.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
						debug = true
					} else {
						rendererb.Destroy()
						windowb.Destroy()
						debug = false
					}
				}
				if player == 0 && t.Type == 769 && t.Keysym.Sym == 'a' {
					capture = [3]int {0,0,0}
					var emptyvalues Board
					values = emptyvalues
					var emptyfreeThrees [2]Board
					freeThrees = emptyfreeThrees
					var emptybetter BoardData
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

		if player_mode == 1 && player == player_two {
			if victory.Todo == true {
				fmt.Printf("IA must play -> x[%d] y [%d]\n", victory.X, victory.Y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", victory.X, victory.Y)
				player = checkRules(&values, &freeThrees, &capture, victory.X, victory.Y, player)
				victory.Todo = false
			} else {
				var x, y int
				x, y, better = search(&values, &freeThrees, player, px, py, 5, &capture)
				fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", x, y)
				if values[y][x] == 0 {
					player = checkRules(&values, &freeThrees, &capture, x, y, player)
				}
			}
			displayAverages()
			resetTimer()
		}

		if player_mode > 0 {
			_ = renderer.SetDrawColor(236, 240, 241, 0)
			renderer.Clear()
			drawGrid(renderer)
			drawClic(renderer, &values, &capture, &freeThrees)
			if player == 0 {
				drawRestartPanel(renderer)
			}
			renderer.Present()
			if debug == true {
				_ = rendererb.SetDrawColor(236, 240, 241, 0)
				rendererb.Clear()
				drawGrid(rendererb)
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
