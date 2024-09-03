package main

import (
	"fmt"
	"sync"
)

var ps []p
var fs []f
var wg sync.WaitGroup

type p struct {
	id  int
	ch  chan (f)
	nom int
}

type f struct {
	id int
	ch chan (p)
}

func main() {
	initStructs()
	initThreads()
	wg.Wait()
}

func initStructs() {
	for i := 0; i < 3; i++ {
		ps = append(ps, p{id: i, ch: make(chan f)})
		fs = append(fs, f{id: i, ch: make(chan p)})
	}
}

func initThreads() {
	for _, p := range ps {
		wg.Add(1)
		go initPhiloThread(p)
	}

	for _, f := range fs {
		wg.Add(1)
		go initForkThread(f)
	}
}

func initPhiloThread(p p) {
	defer wg.Done()
	fl := fs[p.id]
	fr := fs[(p.id+1)%len(ps)]

	go checkFork(p, fl)
	go checkFork(p, fr)
}

func initForkThread(f f) {
	defer wg.Done()
	// pl := ps[f.id]
	// pr := ps[(f.id-1)%len(ps)]
}

func checkFork(p p, f f) {
	m := <-f.ch // m is a philosopher
	if m.ch == nil {
		grabFork(p, f)
	}
}

func grabFork(p p, f f) {
	f.ch <- p // push p to f channel
	p.ch <- f // push f to p channel

	if len(p.ch) == 2 {
		eat(p)
	} else {
		think(p)
	}
}

func eat(p p) {
	p.nom++                   // eat
	<-fs[p.id].ch             // leave fl
	<-fs[(p.id+1)%len(ps)].ch // leave fr
	<-p.ch                    // leave 1 fork
	<-p.ch                    // leave other
	fmt.Printf("Philosopher %d ate. Now he is %d full.\n", p.id, p.nom)
}

func think(p p) {
	fmt.Printf("Philosopher %d is thinking.", p.id)
}
