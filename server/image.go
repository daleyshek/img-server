package server

import (
	"image"
	"os"
)

func getImageWidth(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, err
	}
	return c.Width, nil
}
