package utils

import (
	"image"
	"image/jpeg"
	"os"
	"log"
)

func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create file: %s, Error: %v", path, err)
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		log.Printf("Failed to encode image to file: %s, Error: %v", path, err)
		return err
	}

	log.Printf("Image saved successfully at: %s", path)
	return nil
}
