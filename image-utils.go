package main

import "image"

func GetImageColsAndRows(image *image.NRGBA) (int, int) {
	return image.Bounds().Max.X, image.Bounds().Max.Y
}
