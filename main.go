package main

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	empty = 0
	player_one = 1
	player_two = -1
)

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int = 800, 800

func drawGrid(renderer *sdl.Renderer) {

	_ = renderer.SetDrawColor(236, 0, 0, 0)
	for i := 1; i < 20; i++ {
		_ = renderer.DrawLine(40, 40 * i, 40 * 19, 40 * i)
		_ = renderer.DrawLine(40 *i, 40, 40 * i, 40 * 19)
	}
	
}


func drawClic(renderer *sdl.Renderer, values [19][19]int) {

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[j][i] == player_one {
				_ = renderer.SetDrawColor(0, 236, 0, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			} else if values[j][i] == player_two {
				_ = renderer.SetDrawColor(0, 0, 236, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			}
		}
	}	
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

func checkBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func checkRules(values [19][19]int, nb int, y int, x int) ([19][19]int, bool) {

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

	if f(-1, -1) + f(1, 1) >= 4 || f(1, -1) + f(-1, 1) >= 4 ||
		f(0, -1) + f(0, 1) >= 4 || f(-1, 0) + f(1, 0) >= 4 {
		return values, true
	}
	return values, false
}


func checkCaptures(values *[19][19]int, nb int, y int, x int)  {

	fmt.Printf("Played posiiton [%d][%d]\n", x, y)

	forcapture := func (incx, incy int) {
		if !checkBounds(x + 3 * incx, y + 3 * incy) {
			return
		}
		if  values[y + incy][x + incx] == -nb &&
			values[y + 2 * incy][x + 2 * incx] == -nb &&
		 	values[y + 3 * incy][x + 3 * incx] == nb {
				values[y + incy][x + incx] = 0
				values[y + 2 * incy][x + 2 * incx] = 0
		}
	}

	forcapture(-1, -1)
	forcapture(1, 1)
	forcapture(1, -1)
	forcapture(-1, 1)
	forcapture(0, -1)
	forcapture(0, 1)
	forcapture(-1, 0)
	forcapture(1, 0)

	return
}


func gridAnalyse(values [19][19]int, nb int) (int, int) {
	
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
					fmt.Printf("play position x[%d]y[%d] = %d\n", i, j, max)
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


func run() int {

	var event sdl.Event
	var running bool
	var stop bool
	var err error

	var play bool

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
	stop = false
	play = false
	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:

				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)

				if play == true && t.Type == 1025 && stop == false{
					y:= mousePositionToGrid(float64(t.Y))
					x:= mousePositionToGrid(float64(t.X))
					
					if values[y][x] == 0 {
						var rules bool
						values[y][x] = player_one
						fmt.Printf("[%d]\n", values[y][x])
						values, rules = checkRules(values, player_one, y, x)
						checkCaptures(&values, player_one, y, x)
						if rules {
							fmt.Printf("C'est gagne pour le joueur stupide\n")
							stop = true
						}
						play = false
					}
				} else {
					fmt.Printf("Not your turn\n")
				}

			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}

		if (play == false && stop == false) {


			x, y := gridAnalyse(values, player_two)

			if values[y][x] == 0 {
				var rules bool
				values[y][x] = player_two
				fmt.Printf("[%d]\n", values[y][x])
				values, rules = checkRules(values, player_two, y, x)
				checkCaptures(&values, player_two, y, x)
				if rules {
					fmt.Printf("C'est gagne pour l'iA\n")
					stop = true
				}
				play = true
			}
		}

		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		drawGrid(renderer)
		drawClic(renderer, values)
		renderer.Present()
	}

	return 0
}

func main() {
	os.Exit(run())
}