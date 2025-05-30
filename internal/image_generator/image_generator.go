package image_generator

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

type lightState struct {
	Lights [3]bool
	Arrow  bool
}

var tl1SectionStates = [3][3]bool{
	{true, false, false},
	{false, true, false},
	{false, false, true},
}

func TrafficLight1Image(nextState int) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 60))
	bg := image.NewUniform(color.RGBA{255, 255, 255, 255})
	draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

	colors := []color.RGBA{
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
	}

	for i := range 3 {
		x := 10
		y := i*20 + 10
		r := 10

		isActive := tl1SectionStates[nextState-1][i]

		fillColor := color.RGBA{255, 255, 255, 255}
		if isActive {
			fillColor = colors[i]
		}

		drawCircle(img, x, y, r, fillColor)
	}

	var buffer bytes.Buffer

	if err := png.Encode(&buffer, img); err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64String, nil
}

var tl2SectionStates = [7]lightState{
	{Lights: [3]bool{true, false, false}, Arrow: false}, // Красный (20с)
	{Lights: [3]bool{true, false, false}, Arrow: true},  // Красный + стрелка (20с)
	{Lights: [3]bool{true, false, false}, Arrow: true},  // Красный + мигающая стрелка (5с)
	{Lights: [3]bool{true, false, false}, Arrow: false}, // Красный (10с)
	{Lights: [3]bool{true, true, false}, Arrow: false},  // Красный + желтый (2с)
	{Lights: [3]bool{false, false, true}, Arrow: false}, // Зеленый (20с)
	{Lights: [3]bool{false, true, false}, Arrow: false}, // Желтый (2с)
}

func TrafficLight2Image(nextState int) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, 40, 60))
	bg := image.NewUniform(color.RGBA{255, 255, 255, 255})
	draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

	colors := []color.RGBA{
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
	}
	arrowColor := color.RGBA{204, 255, 153, 255} // #CCFF99

	state := tl2SectionStates[nextState-1]

	for i := range 3 {
		x, y, r := 10, 10+i*20, 10
		fillColor := color.RGBA{255, 255, 255, 255}
		if state.Lights[i] {
			fillColor = colors[i]
		}
		drawCircle(img, x, y, r, fillColor)
	}

	if state.Arrow {
		drawCircle(img, 30, 50, 10, arrowColor)
	} else {
		drawCircle(img, 30, 50, 10, color.RGBA{255, 255, 255, 255})
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func drawCircle(img *image.RGBA, x, y, r int, fill color.RGBA) {
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			dist := dx*dx + dy*dy
			if dist <= r*r && dist >= (r-1)*(r-1) {
				img.Set(x+dx, y+dy, color.RGBA{0, 0, 0, 255})
			} else if dist < (r-1)*(r-1) {
				img.Set(x+dx, y+dy, fill)
			}
		}
	}
}
