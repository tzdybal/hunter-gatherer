package main

func main() {
	h := newHunters()
	g := newGatherers()

	h.searches <- defaultSearch
	h.gatherers = g

	h.send(huntersNum)
	g.send(gatherersNum)

	waitForParties(h, g)
}
