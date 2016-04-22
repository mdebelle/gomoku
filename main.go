package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
)

const (
	empty = 0
	player_one = 1
	player_two = -1
	searchMaxTime = 500000000 * time.Nanosecond
	searchMaxDepth = 20 
)

const (
	VerticalAxis = 1
	HorizontalAxis = 2
	LeftDiagAxis =	4 // haut droite
	RightDiagAxis =	8 // bas droite
)

// type searchParam struct {
// 	board 		*[19][19]int
// 	stoped		bool
// 	stopTime	time.Time
// }

type mustdo struct {
	Todo	bool
	X, Y    int
}

type Board [19][19]int

var winTitle string = "Go-Gomoku"
var winWidth, winHeight int = 800, 880

var victoir mustdo

func checkBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func checkVictory(values *Board, nb int, y int, x int) bool {
	f := func (incx, incy int) int {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
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

func checkCaptures(values *Board, nb, x, y, incx, incy int) bool {
	checkAxis := func (x, y, incx, incy int) bool {
		if !checkBounds(x - incx, y - incy) || !checkBounds(x + 2 * incx, y + 2 * incy) {
			return false
		}
		if values[y + incy][x + incx] == nb {
			if values[y + 2 * incy][x + 2 * incx] == -nb && values[y - incy][x - incx] == 0 {
				victoir.X = x - incx
				victoir.Y = y - incy
				victoir.Todo = true
				return true
			} else if values[y + 2 * incy][x + 2 * incx] == 0   && values[y - incy][x - incx] == -nb {
				victoir.X = x + 2 * incx
				victoir.Y = y + 2 * incy
				victoir.Todo = true
				return true
			}
		}
		return false
	}
	f := func (incx, incy int) bool {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
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

func checkDoubleThree(values, freeThrees *Board, x, y, color int) {

	checkAxis := func(x, y, incx, incy, axis int) {
		checkDirection := func(incx, incy int) (int, int) {
			i, spaces, mine := 1, 0, 0
			for ; i < 5; i++ {
				x, y := x + incx * i, y + incy * i
				if !checkBounds(x, y) || values[y][x] == -color{
					return mine, spaces
				} else if values[y][x] == 0 {
					break
				}
				mine++
			}
			for ; i < 5; i++ {
				x, y := x + incx * i, y + incy * i
				if !checkBounds(x, y) || values[y][x] != 0 {
					break
				}
				spaces++
			}
			return mine, spaces
		}

		if values[y][x] != 0 {
			return
		}
		leftMine, leftSpaces := checkDirection(incx, incy)
		rightMine, rightSpaces := checkDirection(-incx, -incy)
		if leftSpaces == 0 || rightSpaces == 0 {
			return
		} else if leftMine + rightMine == 2  && leftSpaces + rightSpaces > 3 {
			freeThrees[y][x] |= axis
		} else {
			freeThrees[y][x] &= ^axis
		}
	}

	for i := 0; i < 4; i++ {
		checkAxis(x + i, y, 1, 0, HorizontalAxis)
		checkAxis(x - i, y, 1, 0, HorizontalAxis)
		checkAxis(x, y + i, 0, 1, VerticalAxis)
		checkAxis(x, y - i, 0, 1, VerticalAxis)
		checkAxis(x + i, y - i, 1, -1, LeftDiagAxis)
		checkAxis(x - i, y + i, 1, -1, LeftDiagAxis)
		checkAxis(x + i, y + i, 1, 1, RightDiagAxis)
		checkAxis(x - i, y - i, 1, 1, RightDiagAxis)
	}
	return
}

func doCaptures(values *Board, nb int, y int, x int) int {
	forcapture := func (incx, incy int) int {
		if !checkBounds(x + 3 * incx, y + 3 * incy) {
			return 0
		}
		if	values[y + incy][x + incx] == -nb &&
			values[y + 2 * incy][x + 2 * incx] == -nb &&
		 	values[y + 3 * incy][x + 3 * incx] == nb {
				values[y + incy][x + incx] = 0
				values[y + 2 * incy][x + 2 * incx] = 0
				return 2
		}
		return 0
	}
	return  forcapture(-1, -1) + forcapture(1, 1) +
			forcapture(1, -1) + forcapture(-1, 1) +
			forcapture(0, -1) + forcapture(0, 1) +
			forcapture(-1, 0) + forcapture(1, 0)
}

func checkRules(values, freeThrees *Board, capture *[3]int, x, y, player int) int {
	checkDoubleThree(values, freeThrees, x, y, player)
	values[y][x] = player
	victory := checkVictory(values, player, y, x)
	if victory == true {
		fmt.Printf("Victoire \\o/ %d\n", player)
		return 0
	}
	capture[player + 1] += doCaptures(values, player, y, x)
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

func gridAnalyse(values *Board, nb int) (int, int) {
	f := func (incx, incy , x, y, nb int) int {
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
				return i
			}
			x += incx
			y += incy
		}
		return 5
	}
	betterx, bettery := -1, -1
	max := 0
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if (values[i][j] == 0) {
				var t int
				tmp := f(-1, -1, j, i, player_two) + f(1, 1, j, i, player_two)
				t = f(1, -1, j, i, player_two) + f(-1, 1, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = f(0, -1, j, i, player_two) + f(0, 1, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = f(-1, 0, j, i, player_two) + f(1, 0, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = (f(-1, -1, j, i, player_one) + f(1, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(1, -1, j, i, player_one) + f(-1, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(0, -1, j, i, player_one) + f(0, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(-1, 0, j, i, player_one) + f(1, 0, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				if tmp > max {
					max = tmp
					bettery = i
					betterx = j
				}
			}
		}
	}
	if (betterx < 0) {
		betterx = rand.Int() % 19
		bettery = rand.Int() % 19
	}
	return betterx, bettery
}

func init() {
    runtime.LockOSThread()
}

func run() int {
	var event sdl.Event
	var running bool
	var err error
	var player int
	victoir.Todo = false
	var capture [3]int
	var values Board
	var freeThrees Board
	var better [19][19][5]int

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	windowb, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	running = true
	player = 1
	var px, py int
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if player == player_one && t.Type == 1025 {
					py = mousePositionToGrid(float64(t.Y))
					px = mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", px, py)
					if victoir.Todo == true {
						if px == victoir.X && py == victoir.Y {
							player = checkRules(&values, &freeThrees, &capture, px, py, player)
							victoir.Todo = false
						} else {
							fmt.Printf("you must play in [%d][%d]\n", victoir.X, victoir.Y)
						}
					} else if values[py][px] == 0 {
						player = checkRules(&values, &freeThrees, &capture, px, py, player)
					}
				}
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
		if player == player_two {
			if victoir.Todo == true {
				fmt.Printf("IA must play -> x[%d] y [%d]\n", victoir.X, victoir.Y)
				player = checkRules(&values, &freeThrees, &capture, victoir.X, victoir.Y, player)
				victoir.Todo = false	
			} else {
				var x, y int
				x, y, better = search(&values, player, px, py, 4, &capture)
				fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
				if values[y][x] == 0 {
					player = checkRules(&values, &freeThrees, &capture, x, y, player)
				}
			}

		}
		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		drawGrid(renderer)
		drawClic(renderer, &values, &capture)
		renderer.Present()

		_ = rendererb.SetDrawColor(236, 240, 241, 0)
		rendererb.Clear()
		drawGrid(rendererb)
		draweval(rendererb, &better, &freeThrees)
		rendererb.Present()
	}
	return 0
}

func main() {

	os.Exit(run())
}