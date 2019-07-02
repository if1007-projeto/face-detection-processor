package main

import (
	"image"

	pigo "github.com/esimov/pigo/core"
)

var (
	minSize      = 20
	maxSize      = 1000
	shiftFactor  = 0.15
	scaleFactor  = 1.1
	angle        = 0.0
	iouThreshold = 0.2
)

type FaceDetector struct {
	classifier   *pigo.Pigo
	iouThreshold float64
	angle        float64
}

func NewFaceDetector(classifier *pigo.Pigo, iouThreshold float64, angle float64) FaceDetector {
	faceDetector := FaceDetector{classifier, iouThreshold, angle}
	return faceDetector
}

func getDefaultCascadeParams(image *image.NRGBA) pigo.CascadeParams {
	frame := pigo.RgbToGrayscale(image)
	cols, rows := GetImageColsAndRows(image)

	cParams := pigo.CascadeParams{
		MinSize:     minSize,
		MaxSize:     maxSize,
		ShiftFactor: shiftFactor,
		ScaleFactor: scaleFactor,
		ImageParams: pigo.ImageParams{
			Pixels: frame,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}
	return cParams
}

func (fd FaceDetector) GetAllFacesPositions(image *image.NRGBA) []pigo.Detection {
	cParams := getDefaultCascadeParams(image)
	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := fd.classifier.RunCascade(cParams, fd.angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = fd.classifier.ClusterDetections(dets, fd.iouThreshold)

	return dets
}
