package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// location of well known files, from 24-ECIPURI
const (
	SpecsFile      = ".well-known/ecips/specs.json"
	RegistriesFile = ".well-known/ecips/known.json"
)

// Spec is the meta data of Specification from 24-ECIPURI
type Spec struct {
	URI           string `json:uri`           // The URI for the specification
	DocumentURL   string `json:documentUrl`   // The location of the authoritative source for the specification.
	DiscussionURL string `json:discussionUrl` // The location of discussion for the specification.
	Status        string `json:status`        // A status description for the specification.
	Author        string `json:author`        // Contact of the author for the specification.
	CreatedAt     string `json:createdAt`     // Date when the specification was created.
}

func main() {
	processRegistry(".")
}

func processRegistry(uri string) {
	var specs []Spec
	var registries []string
	var err error

	err = parseFile(SpecsFile, &specs)
	if err != nil {
		log.Fatal(err)
	}

	err = parseFile(RegistriesFile, &registries)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Registry:", uri)
	log.Println("  specs:", len(specs))
	log.Println("  registries:", len(registries))
}

func parseFile(fileName string, object interface{}) error {
	jsonData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, object)
	if err != nil {
		return err
	}

	return nil
}
