package main

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/veandco/go-sdl2/sdl"
	"os"
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
			if values[j][i] == 1{
				_ = renderer.SetDrawColor(0, 236, 0, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			} else if values[j][i] == 2 {
				_ = renderer.SetDrawColor(0, 0, 236, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			}
		}
	}	
}

func mousePositionToGrid(val float64) int {
	

	fmt.Printf("What is val: %f\n", val)
	t := int(math.Floor((val - 20) / 40))
	fmt.Printf("What is t: %d\n", t)
	if t < 0 {
		t = 0
	} else if t > 18{
		t = 18
	}
	return t
}


func checkRules(values [19][19]int , nb int, y int, x int) ([19][19]int, int) {

	success := true

	var limit int

	if x - 5 < 0 || y - 5 < 0 {
		if x > y {
			limit = 0 - y
		} else {
			limit = 0 - x
		}
	} else {
		limit = -5
	}

	fmt.Printf("limit is %d\n", limit)

	for i := limit ; i < (limit + 5); i++ {
		fmt.Printf("x[%d] y[%d]\n", y+i, x+i)
		if values[y+i][x+i] != nb {
			success = false
			break
		}
	}

	if success == true {
		return values, 1
	}
	return values, 0
}


func run() int {

	var event sdl.Event
	var running bool
	var stop bool
	var err error

	var play bool

	values := [19][19]int{{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
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
	play = true
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
						var rules int
						values[y][x] = 1
						fmt.Printf("[%d]\n", values[y][x])
						values, rules = checkRules(values, 1, y, x)
						if rules == 1 {
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
	
			x := rand.Int() % 19
			y := rand.Int() % 19

			if values[y][x] == 0 {
				var rules int
				values[y][x] = 2
				fmt.Printf("[%d]\n", values[y][x])
				values, rules = checkRules(values, 2, y, x)
				if rules == 1 {
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