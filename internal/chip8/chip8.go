package chip8

import (
        "github.com/faiface/pixel"
        "github.com/faiface/pixel/imdraw"
        "github.com/faiface/pixel/pixelgl"
        "golang.org/x/image/colornames"
        "math/rand"
        "os"
        "time"
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
        delay uint8
}

func (chip8 Chip8) dump() bool {
        return true
}

func (chip8 *Chip8) tick() {
        var pc = chip8.cpu.pc
        var o1 = chip8.ram.buf[pc] >> 4
        var o2 = chip8.ram.buf[pc] & 0xf
        var o3 = chip8.ram.buf[pc+1] >> 4
        var o4 = chip8.ram.buf[pc+1] & 0xf

        _opcode := uint16(chip8.ram.buf[pc]) << 8 | uint16(chip8.ram.buf[pc + 1])

        logger.Info("opcode:  0x%x", _opcode)
        logger.Info("pc: %d", chip8.cpu.pc)

        // Chip-8 Instruction Set Architecture
        // http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1
        switch o1 {
        case 0x0:
                switch _opcode & 0x000F {
                case 0x0000:
                        logger.Debug("00E0 - CLS")
                        chip8.display.Clear()
                        chip8.cpu.Increment()
                case 0x000E:
                        logger.Debug("00EE - RET")
                        chip8.cpu.sp -= 1
                        chip8.cpu.pc = chip8.cpu.stack[chip8.cpu.sp]
                        chip8.cpu.Increment()
                default:
                        panic("Unknown opcode")
                }

        case 0x1:
                logger.Debug("1nnn - JP addr")
                chip8.cpu.Jump(addr(o2, o3, o4))
        case 0x2:
                logger.Debug("2nnn - CALL addr")
                chip8.cpu.stack[chip8.cpu.sp] = chip8.cpu.pc
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
                switch _opcode & 0x000F {
                case 0x0000:
                        logger.Debug("8xy0 - LD Vx, Vy")
                        chip8.cpu.v[o2] = chip8.cpu.v[o3]
						chip8.cpu.Increment()
                case 0x0001:
                        logger.Debug("8xy1 - OR Vx, Vy")
                        chip8.cpu.v[o2] = chip8.cpu.v[o2] | chip8.cpu.v[o3]
                        chip8.cpu.Increment()
                case 0x0002:
                        logger.Debug("8xy2 - AND Vx, Vy")
                        chip8.cpu.v[o2] = chip8.cpu.v[o2] & chip8.cpu.v[o3]
                        chip8.cpu.Increment()
                case 0x0003:
                        logger.Debug("8xy3 - XOR Vx, Vy")
                        chip8.cpu.v[o2] = chip8.cpu.v[o2] ^ chip8.cpu.v[o3]
                        chip8.cpu.Increment()
                case 0x0004:
                        logger.Debug("8xy4 - ADD Vx, Vy")
                        xy := chip8.cpu.v[o2] + chip8.cpu.v[o3]
                        if xy > 0xff {
                                chip8.cpu.v[0xF] = 1
                        } else {
                                chip8.cpu.v[0xF] = 0
                        }
                        chip8.cpu.v[o2] = xy & 0xff
                        chip8.cpu.Increment()
                case 0x0005:
                        logger.Debug("8xy5 - SUB Vx, Vy")
                        if o3 > o2 {
                                chip8.cpu.v[0xF] = 0
                        } else {
                                chip8.cpu.v[0xF] = 1
                        }
                        chip8.cpu.v[o2] = o2 - o3
                        chip8.cpu.Increment()
                case 0x0006:
                        logger.Debug("8xy6 - SHR Vx {, Vy}")
                        chip8.cpu.v[0xF] = chip8.cpu.v[o2] & 0x1
                        chip8.cpu.v[o2] = chip8.cpu.v[o2] >> 1
                        chip8.cpu.Increment()
                case 0x0007:
                        logger.Debug("8xy7 - SUBN Vx, Vy")
                        if o3 > o2 {
                                chip8.cpu.v[0xF] = 1
                        } else {
                                chip8.cpu.v[0xF] = 0
                        }
                        chip8.cpu.v[o2] = chip8.cpu.v[o3] - chip8.cpu.v[o2]
                        chip8.cpu.Increment()
                case 0x000E:
                        logger.Debug("8xyE - SHL Vx {, Vy}")
                        chip8.cpu.v[0xF] = chip8.cpu.v[o2] >> 7
                        chip8.cpu.v[o2] = chip8.cpu.v[o2] << 1
                        chip8.cpu.Increment()
                default:
                        panic("Unknown opcode")
                }

        case 0x9:
                logger.Debug("9xy0 - SNE Vx, Vy")
                if chip8.cpu.v[o2] != chip8.cpu.v[o3] {
                        chip8.cpu.Skip()
                } else {
                        chip8.cpu.Increment()
                }
        case 0xA:
                logger.Debug("Annn - LD I, addr")
                chip8.cpu.i = addr(o2, o3, o4)
                chip8.cpu.Increment()
        case 0xB:
                logger.Debug("Bnnn - JP V0, addr")
                chip8.cpu.Jump(addr(o2,o3,o4))
        case 0xC:
                logger.Debug("Cxkk - RND Vx, byte")
                rnd := uint8(rand.Intn(256))
                kk := value(o3, o4)
                chip8.cpu.v[o2] = rnd & kk
                chip8.cpu.Increment()
        case 0xD:
                logger.Debug("Dxyn - DRW Vx, Vy, nibble")
                var sprite [256]uint8
                for i := uint8(0); i < uint8(o4); i++ {
                        sprite[i] = chip8.ram.buf[int(chip8.cpu.i) + int(i)]
                }
                vx := chip8.cpu.v[o2]
                vy := chip8.cpu.v[o3]
                collision := chip8.display.Draw(vx, vy, uint8(o4), sprite)
                if collision {
                        chip8.cpu.v[0xF] = 1
                }
                chip8.cpu.Increment()
        case 0xE:
                switch _opcode & 0x00FF {
                case 0x009E:
                        logger.Debug("Ex9E - SKP Vx")
                        vx := chip8.cpu.v[o2]
                        if chip8.cpu.keyboard.keys[vx] == true {
                                chip8.cpu.Skip()
                        } else {
                                chip8.cpu.Increment()
                        }
                case 0x00A1:
                        logger.Debug("ExA1 - SKNP Vx")
                        vx := chip8.cpu.v[o2]
                        if chip8.cpu.keyboard.keys[vx] == false {
                                chip8.cpu.Skip()
                        } else {
                                chip8.cpu.Increment()
                        }
                default:
                        panic("Unknown opcode")
        }
        case 0xF:
                switch _opcode & 0x00FF {
                case 0x0007:
                        logger.Debug("Fx07 - LD Vx, DT")
                        chip8.cpu.v[o2] = chip8.delay
                        chip8.cpu.Increment()
                case 0x000A:
                        logger.Debug("Fx0A - LD Vx, K")
                        vx := chip8.cpu.v[o2]
                        KeyPressed := false
                        for i:= 0; i < len(chip8.cpu.keyboard.keys); i++ {
                                if chip8.cpu.keyboard.keys[i] == true {
                                        chip8.cpu.v[vx] = uint8(i)
                                        KeyPressed = true
                                }
                        }
                        if !KeyPressed {
                                return
                        }
                        chip8.cpu.Increment()
                case 0x0015:
                        logger.Debug("Fx15 - LD DT, Vx")
                        chip8.delay = chip8.cpu.v[o2]
                        chip8.cpu.Increment()
                case 0x0018:
                        logger.Debug("Fx18 - LD ST, Vx")
                        // FIXME: not implemented
                        chip8.cpu.Increment()
                case 0x001E:
                        logger.Debug("Fx1E - ADD I, Vx")
                        chip8.cpu.i = chip8.cpu.i + uint16(chip8.cpu.v[o2])
                        chip8.cpu.Increment()
                case 0x0029:
                        logger.Debug("Fx29 - LD F, Vx")
                        chip8.cpu.i = uint16(chip8.cpu.v[o2]) * 0x5
                        chip8.cpu.Increment()
                case 0x0033:
                        logger.Debug("Fx33 - LD B, Vx")
                        vx := chip8.cpu.v[o2]
                        chip8.ram.buf[chip8.cpu.i] = (vx / 100) % 10
                        chip8.ram.buf[chip8.cpu.i+1] = (vx / 10) % 10
                        chip8.ram.buf[chip8.cpu.i+2] = vx % 10
                        chip8.cpu.Increment()
                case 0x0055:
                        logger.Debug("Fx55 - LD [I], Vx")
                        for x := uint8(0); x < o2; x++ {
                                chip8.ram.buf[int(chip8.cpu.i) + int(x)] = chip8.cpu.v[x]
                        }
                        chip8.cpu.Increment()
                case 0x0065:
                        logger.Debug("Fx65 - LD Vx, [I]")
                        for x := uint8(0); x < o2; x++ {
                                chip8.cpu.v[x] = chip8.ram.buf[int(chip8.cpu.i) + int(x)]
                        }
                        chip8.cpu.Increment()
                default:
                        panic("Unknown opcode")
                }
        default:
                panic("Unknown opcode")
        }

        if chip8.delay > 0 {
                chip8.delay -= 1
        }
}

func (chip8 *Chip8) peepEvent(win *pixelgl.Window) {
        if win.JustPressed(pixelgl.KeyEscape) {
                win.SetClosed(true)
        }

        switch true {
        case win.JustPressed(pixelgl.Key1):
                chip8.cpu.keyboard.KeyDown(0x0)
        case win.JustPressed(pixelgl.Key2):
                chip8.cpu.keyboard.KeyDown(0x1)
        case win.JustPressed(pixelgl.Key3):
                chip8.cpu.keyboard.KeyDown(0x2)
        case win.JustPressed(pixelgl.Key4):
                chip8.cpu.keyboard.KeyDown(0x3)
        case win.JustPressed(pixelgl.KeyQ):
                chip8.cpu.keyboard.KeyDown(0x4)
        case win.JustPressed(pixelgl.KeyW):
                chip8.cpu.keyboard.KeyDown(0x5)
        case win.JustPressed(pixelgl.KeyE):
                chip8.cpu.keyboard.KeyDown(0x6)
        case win.JustPressed(pixelgl.KeyR):
                chip8.cpu.keyboard.KeyDown(0x7)
        case win.JustPressed(pixelgl.KeyA):
                chip8.cpu.keyboard.KeyDown(0x8)
        case win.JustPressed(pixelgl.KeyS):
                chip8.cpu.keyboard.KeyDown(0x9)
        case win.JustPressed(pixelgl.KeyD):
                chip8.cpu.keyboard.KeyDown(0xA)
        case win.JustPressed(pixelgl.KeyF):
                chip8.cpu.keyboard.KeyDown(0xB)
        case win.JustPressed(pixelgl.KeyZ):
                chip8.cpu.keyboard.KeyDown(0xC)
        case win.JustPressed(pixelgl.KeyX):
                chip8.cpu.keyboard.KeyDown(0xD)
        case win.JustPressed(pixelgl.KeyC):
                chip8.cpu.keyboard.KeyDown(0xE)
        case win.JustPressed(pixelgl.KeyV):
                chip8.cpu.keyboard.KeyDown(0xF)
        }

        switch true {
        case win.JustPressed(pixelgl.Key1):
                chip8.cpu.keyboard.KeyUp(0x0)
        case win.JustPressed(pixelgl.Key2):
                chip8.cpu.keyboard.KeyUp(0x1)
        case win.JustPressed(pixelgl.Key3):
                chip8.cpu.keyboard.KeyUp(0x2)
        case win.JustPressed(pixelgl.Key4):
                chip8.cpu.keyboard.KeyUp(0x3)
        case win.JustPressed(pixelgl.KeyQ):
                chip8.cpu.keyboard.KeyUp(0x4)
        case win.JustPressed(pixelgl.KeyW):
                chip8.cpu.keyboard.KeyUp(0x5)
        case win.JustPressed(pixelgl.KeyE):
                chip8.cpu.keyboard.KeyUp(0x6)
        case win.JustPressed(pixelgl.KeyR):
                chip8.cpu.keyboard.KeyUp(0x7)
        case win.JustPressed(pixelgl.KeyA):
                chip8.cpu.keyboard.KeyUp(0x8)
        case win.JustPressed(pixelgl.KeyS):
                chip8.cpu.keyboard.KeyUp(0x9)
        case win.JustPressed(pixelgl.KeyD):
                chip8.cpu.keyboard.KeyUp(0xA)
        case win.JustPressed(pixelgl.KeyF):
                chip8.cpu.keyboard.KeyUp(0xB)
        case win.JustPressed(pixelgl.KeyZ):
                chip8.cpu.keyboard.KeyUp(0xC)
        case win.JustPressed(pixelgl.KeyX):
                chip8.cpu.keyboard.KeyUp(0xD)
        case win.JustPressed(pixelgl.KeyC):
                chip8.cpu.keyboard.KeyUp(0xE)
        case win.JustPressed(pixelgl.KeyV):
                chip8.cpu.keyboard.KeyUp(0xF)
        }
}

func handle(chip8 Chip8, win *pixelgl.Window) {
        for !win.Closed() {
                chip8.tick()
                display := chip8.display.Render()

                win.Clear(colornames.Black)

                for x := 0; x < 64; x++ {
                        for y := 0; y < 32; y++ {
                                if !display[x][y] {
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
                chip8.peepEvent(win)
                win.Update()

                // 60Hz
                time.Sleep((1000 / 120) * time.Millisecond)
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

        file, err := os.ReadFile("../../roms/PONG")

        chip8 := Chip8{}

        // FIXME: initialize
        chip8.cpu.pc = uint16(0x200)

        chip8.ram.Load(file)

        handle(chip8, win)
}
