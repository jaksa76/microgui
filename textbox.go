package microgui

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// TextBox is a multiline Label. It will word-wrap the given content.
// The conent will never cross the boundaries of the TextBox, and will be truncated if longer.
type TextBox struct {
	text   string
	bounds image.Rectangle
}

// NewTextBox creates a new TextBox
func NewTextBox(text string, x, y, w, h int) *TextBox {
	return &TextBox{text, image.Rect(x, y, x+w, y+h)}
}

func (t *TextBox) handleInput(input *userInput) {}

func (t *TextBox) draw(img *ebiten.Image) {
	drawBorder(img, t.bounds)
	textX := t.bounds.Min.X + 2
	for i := 0; i < t.heightInChars() && len(t.getLine(i)) != 0; i++ {
		textY := t.bounds.Min.Y + i*12
		ebitenutil.DebugPrintAt(img, t.getLine(i), textX, textY)
	}
}

func (t *TextBox) getLine(n int) string {
	words := strings.Split(t.text, " ")
	lineNo, wordNo := 0, 0
	var currentLineText string
	for lineNo <= n {
		currentLineText = ""
		for wordNo < len(words) && len(currentLineText)+len(words[wordNo])+1 < t.widthInChars() {
			if len(currentLineText) != 0 {
				currentLineText += " "
			}
			currentLineText += words[wordNo]
			wordNo++
		}
		lineNo++
	}
	return currentLineText
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

// SetWidth asynchronously sets the width of the TextBox
func (t *TextBox) SetWidth(w int) {
	updates <- func() { t.bounds.Max.X = t.bounds.Min.X + w }
}

// SetHeight asynchronously sets the height of the TextBox
func (t *TextBox) SetHeight(h int) {
	updates <- func() { t.bounds.Max.Y = t.bounds.Min.Y + h }
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
