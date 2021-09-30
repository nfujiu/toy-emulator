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

func addr(n1 uint8, n2 uint8, n3 uint8) uint16 {
        return uint16(n1) << 8 + uint16(n2) << 4 + uint16(n3)
}

type Chip8 struct {
        cpu Cpu
        ram Ram
        display Display
}

func (chip8 Chip8) dump() bool {
        return true
}

func (chip8 *Chip8) tick() {
        var pc uint16 =  uint16(0x200) + chip8.cpu.pc

        var o1 uint8 = chip8.ram.buf[pc] >> 4
        var o2 uint8 = chip8.ram.buf[pc] >> 0xf
        var o3 uint8 = chip8.ram.buf[pc+1] >> 4
        var o4 uint8 = chip8.ram.buf[pc+1] >> 0xf
        println(o1)
        println(o2)
        println(o3)
        println(o4)

        // Chip-8 Instruction Set Architecture
        // http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1
        switch o1 {
        case 0x0:
                println(o1, o2, o3, o4)
        case 0x1:
                chip8.cpu.pc = addr(o2, o3, o4)
        case 0x2:
                chip8.cpu.stack[chip8.cpu.sp] = pc
                chip8.cpu.sp += 1
                chip8.cpu.pc = addr(o2, o3, o4)
        case 0x3:
                println(o1, o2, o3, o4)
        case 0x4:
                println(o1, o2, o3, o4)
        case 0x5:
                println(o1, o2, o3, o4)
        case 0x6:
                println(o1, o2, o3, o4)
        case 0x7:
                println(o1, o2, o3, o4)
        case 0x8:
                println(o1, o2, o3, o4)
        case 0x9:
                println(o1, o2, o3, o4)
        case 0xA:
                println(o1, o2, o3, o4)
        case 0xB:
                println(o1, o2, o3, o4)
        case 0xC:
                println(o1, o2, o3, o4)
        case 0xD:
                println(o1, o2, o3, o4)
        case 0xE:
                println(o1, o2, o3, o4)
        case 0xF:
                println(o1, o2, o3, o4)
        default:
                logger.Debug("Unknown opcode")
        }
        // opcode := uint16(chip8.ram.buf[pc]) << 8 | uint16(chip8.ram.buf[pc + 1])
        // decode := opcode & 0xF000

        // switch opcode & 0xF000 {
        // case 0x6000:
        //    println(decode)
        // default:
        //    println("Unknown opcode: 0x%X\n", opcode)
        // }

        // Chip-8 Instruction Set Architecture
        // http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1

        chip8.cpu.Increment()
}

func handle(chip8 Chip8, win *pixelgl.Window) {
        for !win.Closed() {
                chip8.tick()

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

func Run() {
        cfg := pixelgl.WindowConfig{
                Title:  "chip8 Emulator",
                Bounds: pixel.R(0, 0, windowWidth, windowHeight),
                VSync: true,
        }
        win, err := pixelgl.NewWindow(cfg)
        if err != nil {
                panic(err)
        }

        file, err := ioutil.ReadFile("../../roms/PONG")

        chip8 := Chip8{}
        chip8.ram.Load(file)

        handle(chip8, win)
}
