package microgui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// TextField displays a string and lets the user modify it
type TextField struct {
	bounds         image.Rectangle
	content        string
	hasFocus       bool
	cursorPosition int
	contentOffset  int
}

// NewTextField creates a text field
func NewTextField(content string, x, y, w, h int) *TextField {
	return &TextField{image.Rect(x, y, x+w, y+h), content, false, 0, 0}
}

// SetContent asynchronously sets the content of the text field
func (t *TextField) SetContent(s string) {
	updates <- func() { t.content = s }
}

func (t *TextField) handleInput(input *userInput) {
	if input.lmb.clicked {
		t.hasFocus = input.mousePosition.In(t.bounds)
	}

	if t.hasFocus {
		if repeatingKeyPressed(ebiten.KeyBackspace) && t.content != "" {
			if t.cursorPosition > 0 {
				t.content = t.content[:t.cursorPosition-1] + t.content[t.cursorPosition:]
				t.cursorPosition--
			}
		} else if repeatingKeyPressed(ebiten.KeyDelete) && t.content != "" {
			if t.cursorPosition < len(t.content) {
				t.content = t.content[:t.cursorPosition] + t.content[t.cursorPosition+1:]
			}
		} else if repeatingKeyPressed(ebiten.KeyLeft) {
			if t.cursorPosition > 0 {
				t.cursorPosition--
			}
		} else if repeatingKeyPressed(ebiten.KeyRight) {
			if t.cursorPosition < len(t.content) {
				t.cursorPosition++
			}
		} else if repeatingKeyPressed(ebiten.KeyHome) {
			t.cursorPosition = 0
		} else if repeatingKeyPressed(ebiten.KeyEnd) {
			t.cursorPosition = len(t.content)
		} else {
			added := ebiten.InputChars()
			t.content = t.content[:t.cursorPosition] + string(added) + t.content[t.cursorPosition:]
			t.cursorPosition += len(added)
		}

		// scroll to cursor
		if t.cursorPosition < t.contentOffset {
			t.contentOffset = t.cursorPosition
		} else if t.cursorPosition > t.contentOffset+t.widthInChars() {
			t.contentOffset = t.cursorPosition - t.widthInChars()
		}
	}
}

func (t *TextField) draw(img *ebiten.Image) {
	b := t.bounds
	if !b.In(img.Bounds()) {
		return
	}

	col := color.RGBA{20, 20, 20, 0xff}
	ebitenutil.DrawRect(img, float64(b.Min.X), float64(b.Min.Y), float64(b.Dx()), float64(b.Dy()), col)

	ebitenutil.DebugPrintAt(img, truncToPixels(t.content[t.contentOffset:], b.Dx()), b.Min.X+3, b.Min.Y+b.Dy()/2-8)

	if t.hasFocus {
		cursorX := float64(b.Min.X + 5 + (t.cursorPosition-t.contentOffset)*6)
		cursorY := float64(b.Min.Y + b.Dy()/2 - 6)
		ebitenutil.DrawLine(img, cursorX, cursorY, cursorX, cursorY+12.0, color.White)
	}

	x1 := float64(b.Min.X)
	y1 := float64(b.Min.Y)
	x2 := float64(b.Max.X)
	y2 := float64(b.Max.Y)
	ebitenutil.DrawLine(img, x1, y1, x2, y1, boderDarkColor)       // top
	ebitenutil.DrawLine(img, x1, y2-1, x2, y2-1, borderLightColor) // bottom
	ebitenutil.DrawLine(img, x1+1, y1, x1+1, y2, boderDarkColor)   // left
	ebitenutil.DrawLine(img, x2, y1, x2, y2, borderLightColor)     // right
}

func (t *TextField) widthInChars() int {
	return t.bounds.Dx()/6 - 1
}

func truncToPixels(s string, pixels int) string {
	if len(s) < pixels/6-1 {
		return s
	}
	return s[:pixels/6-1]
}

func (t *TextField) clicked(x, y int)    {}
func (t *TextField) mouseReleased()      {}
func (t *TextField) keyPressed()         {}
func (t *TextField) mouseMoved(x, y int) {}
