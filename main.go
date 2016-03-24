package main

import (
	"fmt"
	"math"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int = 800, 800

func drawGrid(renderer *sdl.Renderer) {

	for i := 1; i < 20; i++ {
		_ = renderer.DrawLine(40, 40 * i, 40 * 19, 40 * i)
		_ = renderer.DrawLine(40 *i, 40, 40 * i, 40 * 19)
	}
	
}


func drawClic(renderer *sdl.Renderer, values [19][19]int) {


	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[i][j] == 1{
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			}
		}
	}

	
}


func run() int {

	var event sdl.Event
	var running bool
	var err error

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
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:

				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				y:=int(math.Floor(float64(t.Y) / 40) - 1)
				x:=int(math.Floor(float64(t.X) / 40) - 1)
				fmt.Printf("[%d][%d]\n", x, y)
				values[y][x] = 1
				fmt.Printf("[%d]\n", values[y][x])

			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		_ = renderer.SetDrawColor(236, 0, 0, 0)
		drawGrid(renderer)
		_ = renderer.SetDrawColor(0, 236, 0, 0)
		drawClic(renderer, values)
		renderer.Present()
	}

	return 0
}

func main() {
	os.Exit(run())
}