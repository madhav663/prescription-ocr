package utils

import (
	"image"
	"image/jpeg"
	"os"
)


func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}