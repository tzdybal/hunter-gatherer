package main

import (
	"sync"
)

const (
	caveDir       = "archive"
	huntersNum    = 10
	gatherersNum  = 50
	defaultSearch = "."
)

type party interface {
	send(count uint)
	wait()
}

func waitForParties(parties ...party) {
	wg := new(sync.WaitGroup)
	wg.Add(len(parties))
	for _, p := range parties {
		go func(p party) {
			defer wg.Done()
			p.wait()
		}(p)
	}
	wg.Wait()
}
