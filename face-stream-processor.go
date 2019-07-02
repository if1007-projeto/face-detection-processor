package main

import (
	"bytes"
	"context"
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

type FaceStreamProcessor struct {
	faceDetector  FaceDetector
	brokers       []string
	producerTopic goka.Stream
	jpgQuality    int
	emitter       *goka.Emitter
}

func NewFaceStreamProcessor(faceDetector FaceDetector, brokers []string, producerTopic goka.Stream, jpgQuality int) FaceStreamProcessor {
	emitter, err := goka.NewEmitter(brokers, producerTopic, new(codec.Bytes))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}

	faceProcessor := FaceStreamProcessor{faceDetector, brokers, producerTopic, jpgQuality, emitter}

	return faceProcessor
}

// emits a single message and leave
func (fp FaceStreamProcessor) publishFace(face image.Image) error {
	var opt jpeg.Options
	opt.Quality = fp.jpgQuality
	faceBuffer := new(bytes.Buffer)

	err := jpeg.Encode(faceBuffer, face, &opt)
	if err != nil {
		log.Printf("[ERROR] encoding image %v", err)
		return err
	}

	err = fp.emitter.EmitSync("", faceBuffer.Bytes())
	if err != nil {
		log.Printf("[ERROR] publishing image: %v", err)
		return err
	}

	return nil
}

func (fp FaceStreamProcessor) publishAllFaces(faces []image.Image) {
	for i := 0; i < len(faces); i++ {
		face := faces[i]
		fp.publishFace(face)
	}
}

// process callback is invoked for each message delivered from
// consumer topic
func (fp FaceStreamProcessor) processFrame(ctx goka.Context, msg interface{}) {
	imageBytes := msg.([]byte)

	image, err := jpeg.Decode(bytes.NewReader(imageBytes))
	sourceImage := pigo.ImgToNRGBA(image)

	if err != nil {
		log.Printf("[ERROR] consuming frame from kafka. %v", err)
	}

	dets := fp.faceDetector.GetAllFacesPositions(sourceImage)

	faces := CropAllFaces(sourceImage, dets)

	if len(faces) > 0 {
		log.Printf("[INFO] %d faces detected", len(faces))
	}

	fp.publishAllFaces(faces)
}

// process messages until ctrl-c is pressed
func (fp FaceStreamProcessor) RunProcessor(consumerTopic goka.Stream, consumerGroup goka.Group) {
	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(consumerGroup,
		goka.Input(consumerTopic, new(codec.Bytes), fp.processFrame),
		goka.Persist(new(codec.Int64)),
	)

	p, err := goka.NewProcessor(fp.brokers, g)
	if err != nil {
		log.Fatalf("[ERROR] creating processor: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Fatalf("[ERROR] running processor: %v", err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait   // wait for SIGINT/SIGTERM
	cancel() // gracefully stop processor
	<-done
}
