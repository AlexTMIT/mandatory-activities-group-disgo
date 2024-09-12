package main

import (
	"fmt"
	"sync"
)

var ps []p // array of philosopher
var fs []f // array of forks
var wg sync.WaitGroup

type p struct { // philosopher
	id  int
	ch  chan f
	nom int
}

type f struct { // fork
	id int
	ch chan p
}

func main() {
	initStructs()
	initThreads()
	wg.Wait()
}

func initStructs() {
	// initalises and adds p's and f's to array
	for i := 1; i <= 5; i++ {
		ps = append(ps, p{id: i, ch: make(chan f, 2)}) // 1 p needs 2 forks
		fs = append(fs, f{id: i, ch: make(chan p, 1)}) // only 1 f is handed to 1 p
	}
}

func initThreads() {
	for _, p := range ps {
		wg.Add(1)
		fmt.Printf("*** Starting thread %d\n", p.id)
		go run(&p)
	}
}

func run(p *p) {
	defer wg.Done()
	for p.nom < 3 {
		//fmt.Printf("* Starting iteration of philosopher %d with nom %d\n", p.id, p.nom)

		fl := fs[p.id-1]
		fr := fs[(p.id)%len(ps)]

		//fmt.Printf("Philosopher %d is checking left fork%d\n", p.id, fr.id)
		checkFork(p, fl)
		//fmt.Printf("Philosopher %d is checking right fork%d\n", p.id, fr.id)
		checkFork(p, fr)
	}
	fmt.Printf("Philosopher %d is done eating.\n", p.id)
}

func checkFork(p *p, f f) {
	if len(f.ch) == 0 {
		fmt.Printf("Fork %d has no philosopher, so philosopher %d attempts to grab it.\n", f.id, p.id)
		grabFork(p, f)
	}
}

func grabFork(p *p, f f) {
	f.ch <- *p // push p to f channel
	p.ch <- f  // push f to p channel

	if len(p.ch) == 2 {
		eat(p)
	} else {
		think(p)
	}
}

func eat(p *p) {
	p.nom++                 // eat
	<-fs[p.id-1].ch         // leave fl
	<-fs[(p.id)%len(ps)].ch // leave fr
	<-p.ch                  // 1st fork leaves
	<-p.ch                  // 2nd fork leaves
	fmt.Printf("Philosopher %d ate. Now he is %d full.\n", p.id, p.nom)
}

func think(p *p) {
	fmt.Printf("Philosopher %d is thinking.\n", p.id)
}
