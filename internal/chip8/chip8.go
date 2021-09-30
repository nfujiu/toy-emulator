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

func value(n1 uint8, n2 uint8) uint8 {
        return n1 << 4 + n2
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

        // Chip-8 Instruction Set Architecture
        // http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1
        switch o1 {
        case 0x0:

        case 0x1:
                logger.Debug("1nnn - JP addr")
                chip8.cpu.Jump(addr(o2, o3, o4))
        case 0x2:
                logger.Debug("2nnn - CALL addr")
                chip8.cpu.stack[chip8.cpu.sp] = pc
                chip8.cpu.sp += 1
                chip8.cpu.Jump(addr(o2, o3, o4))
        case 0x3:
                logger.Debug("3xkk - SE Vx, byte")
                kk := value(o3, o4)
                vx := chip8.cpu.v[o2]
                if kk == vx {
                        chip8.cpu.Skip()
                } else {
                        chip8.cpu.Increment()
                }
        case 0x4:
                logger.Debug("4xkk - SNE Vx, byte")
                kk := value(o3, o4)
                vx := chip8.cpu.v[o2]
                if kk != vx {
                        chip8.cpu.Skip()
                } else {
                        chip8.cpu.Increment()
                }
        case 0x5:
                logger.Debug("5xy0 - SE Vx, Vy")
                vx := chip8.cpu.v[o2]
                yx := chip8.cpu.v[o3]
                if vx == yx {
                        chip8.cpu.Skip()
                } else {
                        chip8.cpu.Increment()
                }
        case 0x6:
                logger.Debug("6xkk - LD Vx, byte")
                kk := value(o3, o4)
                chip8.cpu.v[o2] = kk
                chip8.cpu.Increment()
        case 0x7:
                logger.Debug("7xkk - ADD Vx, byte")
                kk := value(o3, o4)
                chip8.cpu.v[o2] = chip8.cpu.v[o2] + kk
                chip8.cpu.Increment()
        case 0x8:
                switch o4 {
                case 0x0:
                        logger.Debug("8xy0 - LD Vx, Vy")

                case 0x1:
                case 0x2:
                case 0x3:
                case 0x4:
                case 0x5:
                case 0x6:
                case 0x7:
                case 0xE:
                default:
                        logger.Debug("Unknown opcode")
                }

        case 0x9:

        case 0xA:

        case 0xB:

        case 0xC:

        case 0xD:

        case 0xE:

        case 0xF:

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

        logger.Debug("Current pc: %d", chip8.cpu.pc)
        chip8.cpu.Increment()
}

func (chip8 *Chip8) peepEvent(win *pixelgl.Window) {
        if win.JustPressed(pixelgl.KeyEscape) {
                // TODO: replace panic
                panic(win)
        }

        if win.JustPressed(pixelgl.Key1) {
                logger.Debug("Just Pressed Key1")
        }

        switch true {
        case win.JustPressed(pixelgl.Key1):
                chip8.cpu.Key(0x1)
        case win.JustPressed(pixelgl.Key2):
                chip8.cpu.Key(0x2)
        case win.JustPressed(pixelgl.Key3):
                chip8.cpu.Key(0x3)
        }
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
                chip8.peepEvent(win)
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
