package chip8

import (
        "github.com/faiface/pixel"
        "github.com/faiface/pixel/imdraw"
        "github.com/faiface/pixel/pixelgl"
        "golang.org/x/image/colornames"
)

func Run () {
        cfg := pixelgl.WindowConfig{
                Title:  "chip8 Emulator",
                Bounds: pixel.R(0, 0, 640, 320),
                VSync: true,
        }
        win, err := pixelgl.NewWindow(cfg)
        if err != nil {
                panic(err)
        }

        imd := imdraw.New(nil)
        imd.Color = pixel.RGB(1, 0, 0)
        imd.Push(pixel.V(20, 10))

        imd.Color = pixel.RGB(0, 1, 0)
        imd.Push(pixel.V(80, 10))

        imd.Color = pixel.RGB(0, 0, 1)
        imd.Push(pixel.V(50, 70))

        imd.Rectangle(0)

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
