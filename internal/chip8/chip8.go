package chip8

import (
        "github.com/faiface/pixel"
        "github.com/faiface/pixel/imdraw"
        "github.com/faiface/pixel/pixelgl"
        "golang.org/x/image/colornames"
        "io/ioutil"
)

const (
        windowWidth float64 = 640
        windowHeight float64 = 320
)

type Chip8 struct {
        cpu Cpu
        ram Ram
        display Display
}

func handle (chip8 Chip8, win *pixelgl.Window) {
        for !win.Closed() {
                win.Clear(colornames.Black)

                for x := 0; x < 64; x++ {
                        for y := 0; y < 32; y++ {
                                if x != y {
                                        continue
                                }

                                imd := imdraw.New(nil)
                                imd.Color = colornames.White
                                imd.Push(pixel.V(float64(x*10), float64(320-(y*10))))
                                imd.Push(pixel.V(float64((x+1)*10), float64(320-((y+1)*10))))
                                imd.Rectangle(0)
                                imd.Draw(win)
                        }
                }
                win.Update()
        }
}

func Run () {
        cfg := pixelgl.WindowConfig{
                Title:  "chip8 Emulator",
                Bounds: pixel.R(0, 0, windowWidth, windowHeight),
                VSync: true,
        }
        win, err := pixelgl.NewWindow(cfg)
        if err != nil {
                panic(err)
        }

        file ,err := ioutil.ReadFile("../../roms/PONG")

        chip8 := Chip8{}
        chip8.ram.Load(file)

        handle(chip8, win)
}
