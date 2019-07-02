package main

import (
	"io/ioutil"
	"log"

	pigo "github.com/esimov/pigo/core"
)

func CreateClassifierFromCascadeFile(cascadeFilePath string) *pigo.Pigo {
	cascadeFile, err := ioutil.ReadFile(cascadeFilePath)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	p := pigo.NewPigo()

	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	return classifier
}
