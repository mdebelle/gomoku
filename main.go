package main

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
)

const (
	empty = 0
	player_one = 1
	player_two = -1
	searchMaxTime = 0.5 
	searchMaxDepth = 20 
	
	winPlayer = 1
	capturedByPlayer = 2
	nothing = 3
	capturedByIA = 4
	winIA = 5

)

type searchParam struct {
	board 		*[19][19]int
	stoped		bool
	stopTime	time.Time
}

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int = 800, 880

func checkBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func checkVictory(values *[19][19]int, nb int, y int, x int) bool {
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

func checkCaptures(values *[19][19]int, nb, x, y, incx, incy int) bool {
	checkAxis := func (x, y, incx, incy int) bool {
		if !checkBounds(x - incx, y - incy) || !checkBounds(x + 2 * incx, y + 2 * incy) {
			return false
		}
		return 	values[y + incy][x + incx] == nb &&
			((values[y + 2 * incy][x + 2 * incx] == -nb && values[y - incy][x - incx] == 0) ||
			 (values[y + 2 * incy][x + 2 * incx] == 0   && values[y - incy][x - incx] == -nb))
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

func doCaptures(values *[19][19]int, nb int, y int, x int) int {
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

func checkRules(values *[19][19]int, capture *[3]int, x, y, player int) int {
	values[y][x] = player
	victory := checkVictory(values, player, y, x)
	if victory == true {
		return 0
	}
	capture[player + 1] += doCaptures(values, player, y, x)
	if capture[player + 1] >= 10 {
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

func gridAnalyse(values *[19][19]int, nb int) (int, int) {
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


func search(values *[19][19]int, player int) (int, int) {
	max := 0
	copy := *values
	var x, y int

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if copy[i][j] == 0 {
				copy[i][j] = player
				pts := 0
				if checkVictory(&copy, player, i, j) {
					return j, i
				} else if doCaptures(&copy, player, i, j) > 0 {
					pts = capturedByIA
				}
				p := minimise(&copy, -player)
				if (pts < p) {
					pts = p
				}
				if max < pts {
					fmt.Printf("coordonees [%d][%d] = %d \n", j, i, pts)
					max = pts
					x = j
					y = i
				}
				copy = *values
			}
		}
	}
	return x, y
}

func minimise(values *[19][19]int, player int) int {
	min := 20
	copy := *values

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if copy[i][j] == 0 {
				copy[i][j] = player
				pts := 20
				if checkVictory(&copy, player, i, j) {
					return winPlayer
				} else if doCaptures(&copy, player, i, j) > 0 {
					pts = capturedByPlayer
				}
				p := maximise(&copy, -player)
				if (pts > p) {
					pts = p
				}
				if min > pts {
					min = pts
				}
				copy = *values
			}
		}
	}
	if min == 20 {
		return nothing
	}
	return min
}

func maximise(values *[19][19]int, player int) int {
	max := 0
	copy := *values

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if copy[i][j] == 0 {
				copy[i][j] = player
				pts := 0
				if checkVictory(&copy, player, i, j) {
					fmt.Printf("%v\n%d %d //%d\n", copy, j, i, player)
					return winIA
				} else if doCaptures(&copy, player, i, j) > 0{
					pts = capturedByIA 
				} else {
					pts = nothing
				}
				if max < pts {
					max = pts
				}
				copy = *values
			}
		}
	}
	if max == 0 {
		return nothing
	}
	return max
}


func run() int {
	var event sdl.Event
	var running bool
	var err error
	var player int
	capture := [3]int {0, 0, 0}
	values := [19][19]int { {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
							{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}

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
	running = true
	player = 1
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if player == player_one && t.Type == 1025 {
					y:= mousePositionToGrid(float64(t.Y))
					x:= mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", x, y)
					if values[y][x] == 0 {
						player = checkRules(&values, &capture, x, y, player)
					}
				}
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
		if player == player_two {
			x, y := search(&values, player)
			fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
			if values[y][x] == 0 {
				player = checkRules(&values, &capture, x, y, player)
			}
		}
		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		drawGrid(renderer)
		drawClic(renderer, &values, &capture)
		renderer.Present()
	}
	return 0
}

func main() {
	os.Exit(run())
}