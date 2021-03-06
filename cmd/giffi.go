package main

import (
	"image/gif"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/coocos/giffi/internal"
)

func readFileName() string {
	if len(os.Args) < 2 {
		log.Fatal("Please specify filename")
	}
	return os.Args[1]
}

func loadGif(fileName string) (*gif.GIF, error) {
	var image *gif.GIF

	// Load via HTTP if given file name seems to be a URL
	if _, err := url.ParseRequestURI(fileName); err == nil {
		resp, err := http.Get(fileName)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		image, err = gif.DecodeAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		image, err = gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}
	}

	return image, nil
}

func main() {
	animation, err := loadGif(readFileName())
	if err != nil {
		log.Fatal("Failed to load image:", err)
	}

	terminal, err := internal.ReadTerminalDimensions()
	if err != nil {
		log.Fatal("Failed to read terminal dimensions", err)
	}

	internal.PlayGif(animation, terminal)
}
