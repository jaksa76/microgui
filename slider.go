package microgui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

// Slider shows is a control that can be dragged from left to right. The value ranges from 0 to 100.
type Slider struct {
	bounds         image.Rectangle
	value          int
	handler        func(int)
	isBeingDragged bool
}

// NewSlider creates a new slider.
func NewSlider(x, y, w, h int, handler func(int)) *Slider {
	return &Slider{image.Rect(x, y, x+w, y+h), 100, handler, false}
}

func (s *Slider) handleInput(input *userInput) {
	if input.lmb.clicked && input.mousePosition.In(s.bounds) {
		s.isBeingDragged = true
	}
	if input.lmb.released {
		s.isBeingDragged = false
	}
	if s.isBeingDragged {
		s.value = clamp(0, 100, (input.mousePosition.X-s.bounds.Min.X-3)*100/(s.bounds.Dx()-6))
		go s.handler(s.value)
	}
}

func (s *Slider) draw(screen *ebiten.Image) {
	drawSquare(screen, s.bounds)
	lineX := float64(s.bounds.Min.X + 3 + s.value*(s.bounds.Dx()-6)/100)
	lineY := float64(s.bounds.Min.Y + 3)
	lineY2 := float64(s.bounds.Max.Y - 3)
	ebitenutil.DrawLine(screen, lineX, lineY, lineX, lineY2, color.White)
}

func clamp(min, max, x int) int {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}
