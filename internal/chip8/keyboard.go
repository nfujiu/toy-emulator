package chip8

type Keyboard struct {
		keys [16]bool
}

func (k *Keyboard) Set(n uint8) {
		k.keys[n] = true
}

func (k *Keyboard) KeyUp(n uint8) {
		k.keys[n] = false
}

func (k *Keyboard) KeyDown(n uint8) {
		k.keys[n] = true
}

func (k *Keyboard) Get(n uint8) bool {
		return k.keys[n]
}