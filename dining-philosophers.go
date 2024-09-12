package main

import (
	"fmt"
	"sync"
)

var ps []p // array of philosopher
var fs []f // array of forks
var wg sync.WaitGroup
var amountFinished = 0

type p struct { // philosopher
	id  int
	ch  chan f
	nom int
}

type f struct { // fork
	id    int
	taken bool // is fork taken by phil?
}

func main() {
	initStructs()
	initThreads()
	wg.Wait()
}

func initStructs() {
	for i := 1; i <= 5; i++ {
		ps = append(ps, p{id: i, ch: make(chan f, 2)}) // 1 p needs 2 forks
		fs = append(fs, f{id: i, taken: false})
	}
}

func initThreads() {
	for i := 1; i <= 5; i++ {
		wg.Add(2)
		fmt.Printf("*** Starting thread %d\n", ps[i-1].id)
		go philGo(&ps[i-1])
		go forkGo(&fs[i-1])
	}
}

func philGo(p *p) {
	defer wg.Done()

	for p.nom < 3 {
		checkFork(p)
	}
	amountFinished++

	fmt.Printf("Philosopher %d is done eating.\n", p.id)
}

func forkGo(f *f) {
	defer wg.Done()

	for amountFinished < len(ps) {
		var p1 = ps[f.id-1]
		var p2 = ps[(f.id)%len(ps)]

		if len(p1.ch) < 2 {
			p1.ch <- *f
			f.taken = true

		} else if len(p2.ch) < 2 {
			p2.ch <- *f
			f.taken = true

		} else {
			f.taken = false
		}
	}
}

func checkFork(p *p) {
	think(p)

	if len(p.ch) != 2 {
		return
	}

	<-p.ch
	<-p.ch
	p.nom++
	fmt.Printf("Philosopher %d is now %d full.\n", p.id, p.nom)
}

func think(p *p) {
	fmt.Printf("Philosopher %d is thinking.\n", p.id)
}
