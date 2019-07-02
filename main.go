package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	pigo "github.com/esimov/pigo/core"
)

const banner = `Face Detection Processor`

var (
	// Flags
	source      = flag.String("in", "", "Source image")
	destination = flag.String("out", "", "Destination image")
	cascadeFile = flag.String("cf", "", "Cascade binary file")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(banner))
		flag.PrintDefaults()
	}

	classifier := CreateClassifierFromCascadeFile("./assets/facefinder")
	faceDetector := NewFaceDetector(classifier, iouThreshold, angle)
	faceMarker := NewFaceMarker()

	image, err := pigo.GetImage("image.jpg")
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}

	dets := faceDetector.GetAllFacesPositions(image)

	buff := new(bytes.Buffer)

	faceMarker.drawMarkerJPG(buff, image, dets, 200)

	file, err := os.Create("out.jpg")
	if err != nil {
		log.Fatalf("Error creating file")
	}
	defer file.Close()

	file.Write(buff.Bytes())

}
