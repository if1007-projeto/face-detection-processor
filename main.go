package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lovoo/goka"
)

const banner = `Face Detection Processor`

var (
	// Default values
	defaultConsumerTopic   = "example-stream"
	defaultConsumerGroup   = "example-group"
	defaultProducerTopic   = "faces"
	defaultCascadeFilePath = "./assets/facefinder"
	defaultJpgQuality      = 100
	defaultBrokersURL      = "localhost:9092"
)

var (
	// Flags
	consumerTopicFlag = flag.String("ct", defaultConsumerTopic, "Consumer topic")
	producerTopicFlag = flag.String("pt", defaultProducerTopic, "Producer topic")
	consumerGroupFlag = flag.String("cg", defaultConsumerGroup, "Consumer topic")
	cascadeFilePath   = flag.String("cf", defaultCascadeFilePath, "Cascade binary file")
	jpgQualityFlag    = flag.Int("q", defaultJpgQuality, "JPG quality (0-100)")
	brokers           = flag.String("b", defaultBrokersURL, "Brokers url")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(banner))
		flag.PrintDefaults()
	}

	// Transform all flags
	consumerTopic := goka.Stream(*consumerTopicFlag)
	consumerGroup := goka.Group(*consumerGroupFlag)
	producerTopic := goka.Stream(*producerTopicFlag)
	brokersURL := strings.Split(*brokers, ",")

	classifier := CreateClassifierFromCascadeFile(*cascadeFilePath)
	faceDetector := NewFaceDetector(classifier, iouThreshold, angle)

	faceStreamProcessor := NewFaceStreamProcessor(faceDetector, brokersURL, producerTopic, *jpgQualityFlag)

	faceStreamProcessor.RunProcessor(consumerTopic, consumerGroup)

}
