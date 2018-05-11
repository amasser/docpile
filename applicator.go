package main

import (
	"log"

	"bitbucket.org/jonathanoliver/docpile/generic/applicators"
)

type Applicator struct{}

func SampleApplicator() applicators.Applicator {
	return &Applicator{}
}

func (this *Applicator) Apply(messages []interface{}) {
	for _, message := range messages {
		log.Println("Applying:", message)
	}
}
