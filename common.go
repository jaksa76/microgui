package microgui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 128
)

var (
	borderLightColor color.Color
	boderDarkColor   color.Color
	fgColor          color.Color
	bgColor          color.Color
	darkBgColor      color.Color
	widgets          []Widget
	mouseDown        bool
	updates          chan func()
	input            userInput
)

func init() {
	fgColor = color.RGBA{255, 255, 255, 0xff}
	bgColor = color.RGBA{30, 30, 30, 0xff}
	darkBgColor = color.RGBA{0, 0, 0, 0xff}
	borderLightColor = color.RGBA{0xff, 0xff, 0xff, 50}
	boderDarkColor = color.RGBA{0, 0, 0, 0xff}
	updates = make(chan func(), 1024)
}

// Widget is anyhing that can be added to the UI
type Widget interface {
	handleInput(input *userInput)

	// tells the widget to draw itself on the image
	draw(screen *ebiten.Image)
}

type userInput struct {
	mousePosition image.Point
	lmb           buttonState
}

type buttonState struct {
	down     bool
	clicked  bool
	released bool
}

func drawBorder(img *ebiten.Image, bounds image.Rectangle) {
	if !bounds.Overlaps(img.Bounds()) {
		return
	}

	x1 := float64(bounds.Min.X)
	y1 := float64(bounds.Min.Y)
	x2 := float64(bounds.Max.X)
	y2 := float64(bounds.Max.Y)
	ebitenutil.DrawLine(img, x1, y1, x2, y1, borderLightColor)     // top
	ebitenutil.DrawLine(img, x1, y2-1, x2, y2-1, boderDarkColor)   // bottom
	ebitenutil.DrawLine(img, x1+1, y1, x1+1, y2, borderLightColor) // left
	ebitenutil.DrawLine(img, x2, y1, x2, y2, boderDarkColor)       // right
}

func update(screen *ebiten.Image) error {
	applyUpdates() // these are updates coming from the business domain layer

	gatherUserInput()
	for _, widget := range widgets {
		widget.handleInput(&input)
	}

	if !ebiten.IsDrawingSkipped() {
		screen.Fill(bgColor)
		for _, widget := range widgets {
			widget.draw(screen)
		}

		ebitenutil.DebugPrint(screen, fmt.Sprintf("%0.2f", ebiten.CurrentFPS()))
	}
	return nil
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func gatherUserInput() {
	input.mousePosition.X, input.mousePosition.Y = ebiten.CursorPosition()
	input.lmb.clicked = !input.lmb.down && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	input.lmb.released = input.lmb.down && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	input.lmb.down = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// apply some updates from the update queue
func applyUpdates() {
	for len(updates) > 0 {
		updateFunction := <-updates
		updateFunction()
	}
}

// Add a widget to the UI
func Add(widget Widget) {
	widgets = append(widgets, widget)
}

// OpenWindow opens the main window
func OpenWindow(w, h int, title string) {
	if err := ebiten.Run(update, w, h, 1, title); err != nil {
		panic(err)
	}
}
