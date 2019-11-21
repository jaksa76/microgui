package microgui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Label is just non-interactive text
type Label struct {
	label  string
	bounds image.Rectangle
}

// NewLabel creates a new Label
func NewLabel(label string, x, y, w, h int) *Label {
	return &Label{label, image.Rect(x, y, x+w, y+h)}
}

func (l *Label) handleInput(input *userInput) {}

func (l *Label) draw(img *ebiten.Image) {
	center := l.bounds.Min.Add(l.bounds.Size().Div(2))
	textX := l.bounds.Min.X
	textY := center.Y - 9
	ebitenutil.DebugPrintAt(img, l.label, textX, textY)
}

// SetLable asynchronously sets the label
func (l *Label) SetLabel(s string) {
	updates <- func() { l.label = s }
}
