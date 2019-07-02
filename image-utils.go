package main

import (
	"image"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

func GetImageColsAndRows(image *image.NRGBA) (int, int) {
	return image.Bounds().Max.X, image.Bounds().Max.Y
}

func getFaceRect(detection pigo.Detection) image.Rectangle {
	return image.Rect(
		detection.Col-detection.Scale/2,
		detection.Row-detection.Scale/2,
		detection.Col+detection.Scale/2,
		detection.Row+detection.Scale/2)
}

func CropAllFaces(sourceImage image.Image, detections []pigo.Detection) []image.Image {
	var faces []image.Image

	for i := 0; i < len(detections); i++ {
		detection := detections[i]
		faceRect := getFaceRect(detection)

		face := imaging.Crop(sourceImage, faceRect)

		faces = append(faces, face)
	}

	return faces
}
