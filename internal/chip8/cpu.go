package chip8

type Cpu struct {
		v [16]uint8
		i uint16
		sp uint8
		pc uint16
		stack    [16]uint16
		keyboard Keyboard
}

func (c *Cpu) Increment() {
		c.pc += 2
}

func (c *Cpu) Skip() {
		c.pc += 4
}

func (c *Cpu) Jump(addr uint16) {
		c.pc = addr
}

func (c *Cpu) Tick(ram Ram) {
		println(c)
}