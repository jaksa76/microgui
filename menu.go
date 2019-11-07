package microgui

import (
	"github.com/hajimehoshi/ebiten"
)

// Menu is a collection of buttons
type Menu struct {
	buttons []*Button
}

// NewMenu creates a new Menu
func NewMenu(x, y, w, h int, handler func(string), options ...string) *Menu {
	var buttons []*Button
	buttonH := h / len(options)
	for i, label := range options {
		buttonY := y + i*buttonH
		buttons = append(buttons, NewButton(label, x, buttonY, w, buttonH, createBtnHandler(label, handler)))
	}
	return &Menu{buttons}
}

func (m *Menu) handleInput(input *userInput) {
	for _, btn := range m.buttons {
		btn.handleInput(input)
	}
}

func (m *Menu) draw(img *ebiten.Image) {
	for _, btn := range m.buttons {
		btn.draw(img)
	}
}

func createBtnHandler(option string, handler func(string)) func() {
	return func() {
		handler(option)
	}
}
