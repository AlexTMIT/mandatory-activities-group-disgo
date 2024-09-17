package main

import (
	"fmt"
	"sync"
	"time"
)

var ps []p // philosophers
var fs []f // forks
var amountFinished = 0

var LIMIT = 1000

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
	wg := new(sync.WaitGroup)
	initStructs()
	initThreads(wg)
	wg.Wait()
}

func initStructs() {
	for i := 1; i <= 5; i++ {
		ps = append(ps, p{id: i, ch: make(chan f, 2)})
		fs = append(fs, f{id: i, taken: false})
	}
}

func initThreads(wg *sync.WaitGroup) {
	for i := 1; i <= 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			philGo(&ps[i-1])
		}()
		go func() {
			defer wg.Done()
			forkGo(&fs[i-1])
		}()
	}
}

func philGo(p *p) {
	for p.nom < LIMIT {
		checkFork(p)
		time.Sleep(1000)
	}

	amountFinished++
	fmt.Printf("Phil %d is DONE eating.\n", p.id)

	if amountFinished == len(ps) {
		fmt.Printf("****** ALL PHILOSOPHERS ARE DONE EATING! ******\n")
	}
}

func forkGo(f *f) {
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
	select {
	case p.ch <- *f:
		f.taken = true
	default:
	}
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
	fmt.Printf("Phil %d ate %d times.\n", p.id, p.nom)
}

func think(p *p) {
	fmt.Printf("Phil %d is thinking.\n", p.id)
}
