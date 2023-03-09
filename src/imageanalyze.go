package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func isWide(image_path []byte) (bool, error) {

	reader, err := os.Open(filepath.Join("input", string(image_path)))
	if err != nil {
		return false, err
	}
	defer reader.Close()

	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return false, err
	}

	if float32(img.Width)/float32(img.Height) > 1.5 {
		return true, nil
	}

	return false, nil
}
