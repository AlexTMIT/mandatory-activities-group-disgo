package main

import (
	"fmt"
	"sync"
	"time"
)

var ps []p // philosophers
var fs []f // forks
var wg sync.WaitGroup
var amountFinished = 0

type p struct {
	id  int
	ch  chan f
	nom int
}

type f struct {
	id    int
	taken bool
}

func main() {
	initStructs()
	initThreads()
	wg.Wait()
}

func initStructs() {
	for i := 1; i <= 5; i++ {
		ps = append(ps, p{id: i, ch: make(chan f, 2)})
		fs = append(fs, f{id: i, taken: false})
	}
}

func initThreads() {
	for i := 1; i <= 5; i++ {
		wg.Add(2)
		go philGo(&ps[i-1])
		go forkGo(&fs[i-1])
	}
}

func philGo(p *p) {
	defer wg.Done()

	for p.nom < 1000 {
		checkFork(p)
		time.Sleep(1000)
	}

	amountFinished++
	fmt.Printf("Philosopher %d is done eating.\n", p.id)
}

func forkGo(f *f) {
	defer wg.Done()

	for amountFinished < len(ps) {
		var p1 = getPhilosopher(f.id - 1)
		var p2 = getPhilosopher(f.id)

		if len(p1.ch) < 2 {
			enterChannel(p1, f)
		} else if len(p2.ch) < 2 {
			enterChannel(p2, f)
		} else {
			f.taken = false
		}
	}
}

func getPhilosopher(id int) p {
	return ps[(id)%len(ps)]
}

func enterChannel(p p, f *f) {
	p.ch <- *f
	f.taken = true
}

func checkFork(p *p) {
	think(p)

	if len(p.ch) != 2 {
		return
	}

	eat(p)
}

func eat(p *p) {
	<-p.ch
	<-p.ch
	p.nom++
	fmt.Printf("Philosopher %d is now %d full.\n", p.id, p.nom)
}

func think(p *p) {
	fmt.Printf("Philosopher %d is thinking.\n", p.id)
}
