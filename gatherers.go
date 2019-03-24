package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/tzdybal/hunter-gatherer/ecipuri"
)

type gatherers struct {
	wg   sync.WaitGroup
	stop chan interface{}

	specs chan ecipuri.Spec
}

func newGatherers() *gatherers {
	return &gatherers{
		stop:  make(chan interface{}),
		specs: make(chan ecipuri.Spec, 1),
	}
}

func (g *gatherers) send(count uint) {
	g.stop = make(chan interface{})

	for i := uint(0); i < count; i++ {
		g.wg.Add(1)
		go g.gatherer(i)
	}
}

func (g *gatherers) cancel() {
	close(g.stop)
}

func (g *gatherers) wait() {
	g.wg.Wait()
}

func (g *gatherers) gatherer(n uint) {
	defer g.wg.Done()
	for {
		select {
		case spec := <-g.specs:
			g.processSpec(spec)
			continue
		default:
		}

		select {
		case <-g.stop:
			return
		default:
		}
	}
}

func (g *gatherers) processSpec(spec ecipuri.Spec) {
	fmt.Println("Processing spec:", spec.URI)
	resp, err := http.Get(spec.DocumentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// TODO: add some extra protection
	dir, file := filepath.Split(spec.URI)
	err = os.MkdirAll(filepath.Join(caveDir, dir), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := os.Create(filepath.Join(caveDir, dir, file) + ".html")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := output.Sync()
		if err != nil {
			fmt.Println(err)
		}
		_ = output.Close()
	}()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
}
