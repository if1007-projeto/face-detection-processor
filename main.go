package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

const banner = `Face Detection Processor`

var (
	// Flags
	source      = flag.String("in", "", "Source image")
	destination = flag.String("out", "", "Destination image")
	cascadeFile = flag.String("cf", "", "Cascade binary file")
	jpgQuality  = 100
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(banner))
		flag.PrintDefaults()
	}

	classifier := CreateClassifierFromCascadeFile("./assets/facefinder")
	faceDetector := NewFaceDetector(classifier, iouThreshold, angle)
	// faceMarker := NewFaceMarker()

	sourceImage, err := pigo.GetImage("./assets/image.jpg")
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}

	dets := faceDetector.GetAllFacesPositions(sourceImage)

	// image := faceMarker.drawMarker(sourceImage, dets)

	faces := CropAllFaces(sourceImage, dets)

	file, err := os.Create("out.jpg")
	if err != nil {
		log.Fatalf("Error creating file")
	}
	defer file.Close()

	var opt jpeg.Options
	opt.Quality = jpgQuality

	for i := 0; i < len(faces); i++ {
		imaging.Save(faces[i], strconv.Itoa(i)+".jpg")
	}

	// jpeg.Encode(file, faces[0], &opt)
}
