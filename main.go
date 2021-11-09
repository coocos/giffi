package main

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

const (
	brightnessToAscii = " .:-=+*#%@"
)

type terminalDimensions struct {
	width  uint
	height uint
}

func readFileName() string {
	if len(os.Args) < 2 {
		log.Fatal("Please specify filename")
	}
	return os.Args[1]
}

func loadGif(fileName string) (*gif.GIF, error) {
	var image *gif.GIF

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err = gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func readTerminalDimensions() (terminalDimensions, error) {
	var terminal terminalDimensions
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return terminal, fmt.Errorf("Failed to read terminal dimensions: %w", err)
	}
	dimensions := strings.Split(strings.TrimSuffix(string(out), "\n"), " ")
	if len(dimensions) != 2 {
		return terminal, fmt.Errorf("stty returned an unexpected output: %s", out)
	}

	height, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return terminal, fmt.Errorf("stty returned an unexpected value for height: %s", dimensions[0])
	}
	width, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return terminal, fmt.Errorf("stty returned an unexpected value for width: %s", dimensions[0])
	}

	terminal.width = uint(width)
	terminal.height = uint(height)

	return terminal, nil
}

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
			// FIXME: This is pretty gross
			character := int(math.Floor(brightness / 65025.0 * float64(len(brightnessToAscii))))
			buffer.WriteByte(brightnessToAscii[character%len(brightnessToAscii)])
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

func playGif(animation *gif.GIF, terminal terminalDimensions) {
	// Loop forever because why not
	for {
		for frameNumber, frame := range animation.Image {
			resizedFrame := resizeImage(frame, terminal.width, terminal.height)
			fmt.Print(frameToAscii(resizedFrame))
			time.Sleep(time.Duration(animation.Delay[frameNumber]*10) * time.Millisecond)
		}
	}
}

func main() {
	animation, err := loadGif(readFileName())
	if err != nil {
		log.Fatal("Failed to load image:", err)
	}

	terminal, err := readTerminalDimensions()
	if err != nil {
		log.Fatal("Failed to read terminal dimensions", err)
	}

	playGif(animation, terminal)
}
