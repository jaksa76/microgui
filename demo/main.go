package main

import (
	"fmt"
	"strconv"

	"github.com/jaksa76/microgui"
)

var (
	i  int
	tf *microgui.TextField
)

func main() {
	fmt.Println("starting...")

	microgui.Add(microgui.NewButton("One", 10, 40, 80, 40, handlerFactory("one")))
	microgui.Add(microgui.NewButton("Two", 10, 80, 80, 40, handlerFactory("two")))
	microgui.Add(microgui.NewButton("Three", 10, 120, 80, 40, handlerFactory("three")))
	microgui.Add(microgui.NewTextField("Hajimemashite", 10, 200, 80, 20))

	tf = microgui.NewTextField("", 10, 240, 160, 20)
	microgui.Add(tf)
	microgui.Add(microgui.NewLabel("slide to change text field", 10, 260, 160, 20))
	microgui.Add(microgui.NewSlider(10, 280, 160, 20, sliderHandler))

	microgui.Add(microgui.NewMenu(160, 40, 120, 100, menuHandler, "A", "B", "C"))

	microgui.OpenWindow(640, 480, "Micro GUI Demo")
}

func buttonHandler() {
	microgui.Add(microgui.NewButton(strconv.Itoa(i), 300+(i/10)*20, 20*(i%10), 20, 20, buttonHandler))
	i++
}

func menuHandler(option string) {
	fmt.Println(option + " selected")
}

func handlerFactory(buttonName string) func() {
	return func() {
		fmt.Println(buttonName + " pressed")
		buttonHandler()
	}
}

func sliderHandler(value int) {
	tf.SetContent(strconv.Itoa(value))
}
