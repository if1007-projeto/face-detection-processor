package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/signal"
	"syscall"

	pigo "github.com/esimov/pigo/core"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

const banner = `Face Detection Processor`

var (
	// Flags
	source      = flag.String("in", "", "Source image")
	destination = flag.String("out", "", "Destination image")
	cascadeFile = flag.String("cf", "", "Cascade binary file")
	jpgQuality  = 100
)

var (
	brokers             = []string{"localhost:9092"}
	topic   goka.Stream = "example-stream"
	group   goka.Group  = "example-group"
)

// process messages until ctrl-c is pressed
func runProcessor() {
	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(ctx goka.Context, msg interface{}) {
		log.Printf("reading message from kafka")

		// ctx.Value() gets from the group table the value that is stored for
		// the message's key.
		if val := ctx.Value(); val != nil {
			jpeg.Decode(bytes.NewReader(val.([]byte)))
		}
	}

	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.String), cb),
		goka.Persist(new(codec.Int64)),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait   // wait for SIGINT/SIGTERM
	cancel() // gracefully stop processor
	<-done
}

func publishAllFaces(faces []image.Image) {
	var opt jpeg.Options
	opt.Quality = jpgQuality

	for i := 0; i < len(faces); i++ {
		faceBuffer := new(bytes.Buffer)

		face := faces[i]
		err := jpeg.Encode(faceBuffer, face, &opt)
		if err != nil {
			log.Fatalf("error encoding image")
		}
	}

}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(banner))
		flag.PrintDefaults()
	}

	classifier := CreateClassifierFromCascadeFile("./assets/facefinder")
	faceDetector := NewFaceDetector(classifier, iouThreshold, angle)

	sourceImage, err := pigo.GetImage("./assets/image.jpg")
	if err != nil {
		log.Fatalf("cannot open the image file: %v", err)
	}

	dets := faceDetector.GetAllFacesPositions(sourceImage)

	faces := CropAllFaces(sourceImage, dets)

	publishAllFaces(faces)

	runProcessor() // press ctrl-c to stop

}
