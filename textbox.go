package microgui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type TextBox struct {
	text   string
	bounds image.Rectangle
}

func NewTextBox(text string, x, y, w, h int) *TextBox {
	return &TextBox{text, image.Rect(x, y, x+w, y+h)}
}

func (t *TextBox) handleInput(input *userInput) {}

func (t *TextBox) draw(img *ebiten.Image) {
	textX := t.bounds.Min.X + 2
	for i := 0; i < t.numberOfLines(); i++ {
		textY := t.bounds.Min.Y + i*12
		ebitenutil.DebugPrintAt(img, t.getLine(i), textX, textY)
	}
}

func (t *TextBox) getLine(n int) string {
	start := min(n*t.widthInChars(), len(t.text))
	end := min((n+1)*t.widthInChars(), len(t.text))
	return t.text[start:end]
}

func (t *TextBox) numberOfLines() int {
	if t.widthInChars() < 1 {
		return 0
	}
	return min(len(t.text)/t.widthInChars()+1, t.heightInChars())
}

func (t *TextBox) heightInChars() int {
	return ((t.bounds.Dy() - 6) / 12)
}

func (t *TextBox) widthInChars() int {
	return t.bounds.Dx()/6 - 1
}

// SetText asynchronously sets the text
func (t *TextBox) SetText(s string) {
	updates <- func() { t.text = s }
}

func (t *TextBox) SetWidth(w int) {
	updates <- func() { t.bounds.Max.X = t.bounds.Min.X + w }
}

func (t *TextBox) SetHeight(h int) {
	updates <- func() { t.bounds.Max.Y = t.bounds.Min.Y + h }
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
