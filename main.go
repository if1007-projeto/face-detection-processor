package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"

	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
)

var minSize = 20
var maxSize = 1000
var shiftFactor = 0.15
var scaleFactor = 1.1
var angle = 0.0
var circleMarker = false

var dc *gg.Context

func main() {
	cascadeFile, err := ioutil.ReadFile("facefinder")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	src, err := pigo.GetImage("image.jpg")
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}

	frame := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

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

	p := pigo.NewPigo()

	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)

	dc = gg.NewContext(cols, rows)
	dc.DrawImage(src, 0, 0)

	buff := new(bytes.Buffer)
	drawMarker(dets, buff)
	dc.SavePNG("out.png")

	fmt.Println(dets)
}

// drawMarker mark the detected face region with the provided
// marker (rectangle or circle) and write it to io.Writer.
func drawMarker(detections []pigo.Detection, w io.Writer) {
	// var qThresh float32 = 5.0

	for i := 0; i < len(detections); i++ {
		// if detections[i].Q > qThresh {
		dc.DrawRectangle(
			float64(detections[i].Col-detections[i].Scale/2),
			float64(detections[i].Row-detections[i].Scale/2),
			float64(detections[i].Scale),
			float64(detections[i].Scale),
		)

		dc.SetLineWidth(3.0)
		dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
		dc.Stroke()
		// }
	}
	dc.SavePNG("out.png")
	// return dc.EncodePNG(w)
}
