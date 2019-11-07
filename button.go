package microgui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Button is a clickable area with a label.
type Button struct {
	label   string
	bounds  image.Rectangle
	handler func()
	pressed bool
}

// NewButton creates a new button and returns the pointer
func NewButton(label string, x, y, w, h int, handler func()) *Button {
	return &Button{label, image.Rect(x, y, x+w, y+h), handler, false}
}

func (b *Button) handleInput(input *userInput) {
	if input.lmb.clicked && input.mousePosition.In(b.bounds) {
		b.pressed = true
		go b.handler()
	}
	if input.lmb.released {
		b.pressed = false
	}
}

func (b *Button) draw(img *ebiten.Image) {
	if b.pressed {
		fx := float64(b.bounds.Min.X)
		fy := float64(b.bounds.Min.Y)
		fw := float64(b.bounds.Dx())
		fh := float64(b.bounds.Dy())
		ebitenutil.DrawRect(img, fx, fy, fw, fh, darkBgColor)
	}
	drawSquare(img, b.bounds)
	center := b.bounds.Min.Add(b.bounds.Size().Div(2))
	textX := center.X - len(b.label)*3
	textY := center.Y - 9
	ebitenutil.DebugPrintAt(img, b.label, textX, textY)
}

