package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// well known constants from 24-ECIPURI
const (
	SpecsFile      = ".well-known/ecips/specs.json"
	RegistriesFile = ".well-known/ecips/known.json"
)

// own consts
const ()

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
	var reader func(string) ([]byte, error)
	// TODO: add git and IPFS support
	switch {
	case strings.HasPrefix(uri, "http"): // match both http and https
		reader = func(suffix string) ([]byte, error) {
			resp, err := http.Get(uri + "/" + suffix)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			return ioutil.ReadAll(resp.Body)
		}
	default: // assume file
		reader = func(suffix string) ([]byte, error) {
			path := filepath.Join(uri, suffix)
			return ioutil.ReadFile(path)
		}
	}
	var specs []Spec
	var registries []string
	var err error

	log.Println("Processing registry:", uri)
	err = parseFile(reader, SpecsFile, &specs)
	if err != nil {
		log.Println(err)
	}

	err = parseFile(reader, RegistriesFile, &registries)
	if err != nil {
		log.Println(err)
	}

	log.Println("  specs:", len(specs))
	log.Println("  registries:", len(registries))

	for _, registry := range registries {
		log.Println("->", registry)
		processRegistry(registry)
	}
}

func parseFile(reader func(string) ([]byte, error), fileName string, object interface{}) error {
	jsonData, err := reader(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, object)
	if err != nil {
		return err
	}

	return nil
}
