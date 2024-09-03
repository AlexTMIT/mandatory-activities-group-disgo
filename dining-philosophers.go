package main

var ps []p
var fs []f

type p struct {
	id int
	ch chan (f)
}

type f struct {
	id int
	ch chan (p)
}

func main() {
	initStructs()
	initThreads()
}

func initStructs() {
	for i := 0; i < 3; i++ {
		ps = append(ps, p{id: i, ch: make(chan f)})
		fs = append(fs, f{id: i, ch: make(chan p)})
	}
}

func initThreads() {
	for _, p := range ps {
		go initPhiloThread(p)
	}

	for _, f := range fs {
		go initForkThread(f)
	}
}

func initPhiloThread(p p) {
	isEating := false
	fl := fs[p.id]
	fr := fs[(p.id+1)%len(ps)]
	amount := 0

	//
}

func initForkThread(f f) {
	isUsed := false
	pl := ps[p.id]
	pr := ps[(p.id-1)%len(ps)]

}

// channel size is >= 2, lock
