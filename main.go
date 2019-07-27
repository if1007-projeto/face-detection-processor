package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

const banner = `Face Detection Processor`

var (
	// Default values
	defaultConsumerTopic   = "frames"
	defaultConsumerGroup   = "frames-consumer-group"
	defaultFacesTopic      = "faces"
	defaultMarksTopic      = "marks"
	defaultCascadeFilePath = "./assets/facefinder"
	defaultJpgQuality      = 100
	defaultBrokersURL      = "localhost:9092"
	retries                = 100
	retryTime              = 2 * time.Second
)

var (
	// Flags
	consumerTopicFlag = flag.String("ct", defaultConsumerTopic, "Consumer topic")
	facesTopicFlag    = flag.String("ft", defaultFacesTopic, "Faces topic")
	marksTopicFlag    = flag.String("mt", defaultMarksTopic, "Marks topic")
	consumerGroupFlag = flag.String("cg", defaultConsumerGroup, "Consumer group")
	cascadeFilePath   = flag.String("cf", defaultCascadeFilePath, "Cascade binary file")
	jpgQualityFlag    = flag.Int("q", defaultJpgQuality, "JPG quality (0-100)")
	brokers           = flag.String("b", defaultBrokersURL, "Brokers url")
)

func tryCreateEmitter(brokersURL []string, topic goka.Stream) *goka.Emitter {
	var emitter *goka.Emitter
	for i := 0; i < retries; i++ {
		log.Printf("[INFO] [%d/%d] trying to connect to kafka (%s)", i+1, retries, brokersURL)
		newEmitter, err := goka.NewEmitter(brokersURL, topic, new(codec.Bytes))
		if err == nil {
			emitter = newEmitter
			break
		}
		log.Printf("[ERROR] creating emitter: %v", err)
		time.Sleep(retryTime)
	}
	return emitter
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(banner))
		flag.PrintDefaults()
	}
	flag.Parse()

	// Transform all flags
	consumerTopic := goka.Stream(*consumerTopicFlag)
	consumerGroup := goka.Group(*consumerGroupFlag)
	facesTopic := goka.Stream(*facesTopicFlag)
	marksTopic := goka.Stream(*marksTopicFlag)
	brokersURL := strings.Split(*brokers, ",")

	log.Printf("[INFO] Kafka url: %s", brokersURL)
	log.Printf("[INFO] Reading frames from \"%s\" topic and publishing in \"%s\" and \"%s\"", consumerTopic, facesTopic, marksTopic)

	classifier := CreateClassifierFromCascadeFile(*cascadeFilePath)
	faceDetector := NewFaceDetector(classifier, iouThreshold, angle)

	facesEmitter := tryCreateEmitter(brokersURL, facesTopic)
	marksEmitter := tryCreateEmitter(brokersURL, marksTopic)

	faceStreamProcessor := NewFaceStreamProcessor(faceDetector, brokersURL, facesEmitter, marksEmitter, *jpgQualityFlag)

	faceStreamProcessor.RunProcessor(consumerTopic, consumerGroup)

}
