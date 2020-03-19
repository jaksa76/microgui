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
		if t.cursorPosition < t.contentOffset { // cursor is left of visible text
			t.contentOffset = t.cursorPosition
		} else if t.cursorPosition > t.contentOffset+t.widthInChars() { // cursor is on right of visible text
			t.contentOffset = t.cursorPosition - t.widthInChars()
		} else if t.contentOffset > len(t.content)-t.widthInChars() { // cursor is inside visible text
			t.contentOffset = max(0, len(t.content)-t.widthInChars())
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

	ebitenutil.DebugPrintAt(img, t.visibleText(), b.Min.X+3, b.Min.Y+b.Dy()/2-8)

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

	// ellipses for text out of bounds
	if t.contentOffset > 0 {
		ebitenutil.DrawLine(img, x1+2, y2-3, x1+3, y2-3, fgColor)
		ebitenutil.DrawLine(img, x1+4, y2-3, x1+5, y2-3, fgColor)
		ebitenutil.DrawLine(img, x1+6, y2-3, x1+7, y2-3, fgColor)
	}
	if t.contentOffset < len(t.content)-t.widthInChars() {
		ebitenutil.DrawLine(img, x2-2, y2-3, x2-3, y2-3, fgColor)
		ebitenutil.DrawLine(img, x2-4, y2-3, x2-5, y2-3, fgColor)
		ebitenutil.DrawLine(img, x2-6, y2-3, x2-7, y2-3, fgColor)
	}
}

func (t *TextField) widthInChars() int {
	return t.bounds.Dx()/6 - 1
}

func (t *TextField) visibleText() string {
	if len(t.content)-t.contentOffset < t.widthInChars() {
		return t.content[t.contentOffset:]
	}
	return t.content[t.contentOffset : t.contentOffset+t.widthInChars()]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *TextField) clicked(x, y int)    {}
func (t *TextField) mouseReleased()      {}
func (t *TextField) keyPressed()         {}
func (t *TextField) mouseMoved(x, y int) {}
