package internal

import (
	"fmt"
	"image"
	"image/gif"
	"math"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

const (
	brightnessToAscii = " .:-=+*#%@"
)

func frameToAscii(frame image.Image) string {
	var buffer strings.Builder

	bounds := frame.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// Skip every other row because terminal fonts are usually taller than wide
		if y%2 == 0 {
			continue
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := frame.At(x, y).RGBA()
			brightness := float64((r + g + b) / 3)
			character := int(math.Floor(brightness / 65536 * float64(len(brightnessToAscii))))
			buffer.WriteByte(brightnessToAscii[character])
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func resizeImage(frame image.Image, width uint, height uint) image.Image {
	if height < width {
		return resize.Resize(0, height*2, frame, resize.Lanczos3)
	} else {
		return resize.Resize(width*2, 0, frame, resize.Lanczos3)
	}
}

// PlayGif plays the GIF by mapping pixels to ASCII
func PlayGif(animation *gif.GIF, terminal TerminalDimensions) {

	frameCache := make(map[int]string)

	// Render the loop first to the cache
	for frameNumber, frame := range animation.Image {
		resizedFrame := resizeImage(frame, terminal.width, terminal.height)
		frameCache[frameNumber] = frameToAscii(resizedFrame)
	}

	// Loop forever because why not
	for {
		for frameNumber := range animation.Image {
			fmt.Print(frameCache[frameNumber])
			time.Sleep(time.Duration(animation.Delay[frameNumber]*10) * time.Millisecond)
		}
	}
}
