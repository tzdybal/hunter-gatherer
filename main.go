package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// well known constants from 24-ECIPURI
const (
	SpecsFile      = ".well-known/ecips/specs.json"
	RegistriesFile = ".well-known/ecips/known.json"
)

// own consts
const (
	CaveDir      = "archive"
	HuntersNum   = 5
	GatherersNum = 10
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

// hunter is a channel for registries locations
var hunter chan string

// gaterer is a channel for specs to be collected
var gatherer chan Spec

var stop chan struct{}

var wg sync.WaitGroup

func main() {
	hunter = make(chan string)
	gatherer = make(chan Spec)
	stop = make(chan struct{})

	start()
}

func start() {
	for i := 0; i < GatherersNum; i++ {
		fmt.Println("g", len(gatherer))
		wg.Add(1)
		go func() {
			for {
				select {
				case spec := <-gatherer:
					processSpec(spec)
				case <-stop:
					break
				default:
				}
			}
			wg.Done()
			fmt.Println("Gatherer done", i)
		}()
	}

	for i := 0; i < HuntersNum; i++ {
		fmt.Println("h", len(hunter))
		wg.Add(1)
		go func() {
			for {
				select {
				case uri := <-hunter:
					processRegistry(uri)
				case <-stop:
					break
				default:
				}
			}
			wg.Done()
			fmt.Println("Gatherer done")
		}()
	}

	hunter <- "."

	wg.Wait()
}

func processSpec(spec Spec) {
	// TODO: archive file
	fmt.Println("%+v", spec)
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

	fmt.Println("Processing registry:", uri)
	err = parseFile(reader, SpecsFile, &specs)
	if err != nil {
		fmt.Println(err)
	}

	err = parseFile(reader, RegistriesFile, &registries)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("  specs:", len(specs))
	fmt.Println("  registries:", len(registries))

	for _, spec := range specs {
		gatherer <- spec
	}

	for _, registry := range registries {
		fmt.Println("->", registry)
		hunter <- registry
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
