package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tzdybal/hunter-gatherer/ecipuri"
)

type hunters struct {
	wg       *sync.WaitGroup
	stop     chan interface{}
	searches chan string
	searched sync.Map

	gatherers *gatherers
}

func newHunters() *hunters {
	return &hunters{
		wg:       new(sync.WaitGroup),
		searches: make(chan string, 10),
	}
}

func (h *hunters) send(count uint) {
	h.stop = make(chan interface{})
	for i := uint(0); i < count; i++ {
		h.wg.Add(1)
		go h.hunter(i)
	}
}

func (h *hunters) wait() {
	h.wg.Wait()
	h.gatherers.cancel()
}

func (h *hunters) hunter(n uint) {
	defer h.wg.Done()
	for {
		shouldStop := false
		select {
		case search := <-h.searches:
			shouldStop = !h.processRegistry(search)
		default:
		}

		select {
		case <-h.stop:
			return
		default:
		}

		if shouldStop {
			close(h.stop)
		}
	}
}

func (h *hunters) processRegistry(uri string) bool {
	if _, alreadyDone := h.searched.LoadOrStore(uri, true); alreadyDone {
		return false
	}

	var reader func(string) ([]byte, error)

	switch {
	// TODO: add more sophisticated matching
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		reader = func(suffix string) ([]byte, error) {
			resp, err := http.Get(uri + "/" + suffix)
			if err != nil {
				return nil, err
			}
			defer func() { _ = resp.Body.Close() }()
			return ioutil.ReadAll(resp.Body)
		}
	// TODO: add git and IPFS support
	default: // assume file
		reader = func(suffix string) ([]byte, error) {
			path := filepath.Join(uri, suffix)
			return ioutil.ReadFile(path)
		}
	}
	var specs []ecipuri.Spec
	var registries []string
	var err error

	fmt.Println("Processing registry:", uri)
	err = h.parseFile(reader, ecipuri.SpecsFile, &specs)
	if err != nil {
		fmt.Println(err)
	}

	err = h.parseFile(reader, ecipuri.RegistriesFile, &registries)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("  specs:", len(specs))
	fmt.Println("  registries:", len(registries))

	somethingNew := false
	for _, registry := range registries {
		if _, alreadyDone := h.searched.LoadOrStore(uri, true); alreadyDone {
			somethingNew = true
			h.searches <- registry
		}
	}

	for _, spec := range specs {
		h.gatherers.specs <- spec
	}

	return somethingNew
}

func (h *hunters) parseFile(reader func(string) ([]byte, error), fileName string, object interface{}) error {
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
