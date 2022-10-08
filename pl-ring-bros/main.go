package main

import (
	"flag"
	"log"
	"time"
)

type Token struct {
	Data     string
	Reciever uint
	TTL      uint
}

type Node struct {
	id      uint
	beforeC <-chan Token
	nextC   chan Token

	doneC <-chan bool
}

func (node *Node) Run() {
	go func() {
		for {
			select {
			case t := <-node.beforeC:
				switch {
				case t.Reciever == node.id:
					log.Printf("Token has been accepted by %d; message: %s (with left ttl = %d)", t.Reciever, t.Data, t.TTL)
				case t.TTL > 0:
					t.TTL -= 1
					node.nextC <- t
				}
			case <-node.doneC:
				return
			}
		}
	}()
}

type TokenRing struct {
	doneC <-chan bool
	nodes []*Node
}

func NewTokenRing(size uint, doneC <-chan bool) *TokenRing {
	tr := &TokenRing{
		doneC: doneC,
		nodes: make([]*Node, 0, size),
	}

	if size < 2 {
		log.Fatal("Expected to have at least 2 size of ring")
	}

	firstNode := &Node{
		id:    0,
		nextC: make(chan Token),
		doneC: doneC,
	}
	tr.nodes = append(tr.nodes, firstNode)

	var i uint
	for i = 1; i < size; i++ {
		tr.nodes = append(tr.nodes, &Node{
			id:      i,
			beforeC: tr.nodes[i-1].nextC,
			nextC:   make(chan Token),
			doneC:   doneC,
		})
	}
	// Assign firstNode beforeC from last node
	firstNode.beforeC = tr.nodes[size-1].nextC

	return tr
}

func (tr *TokenRing) Run() chan Token {
	for _, node := range tr.nodes {
		node.Run()
	}
	return tr.nodes[len(tr.nodes)/2].nextC
}

func main() {
	var chainSize = flag.Uint("n", 12, "Chain size")
	flag.Parse()

	if chainSize == nil {
		log.Fatal("Expected to provide n uint value, number of nodes in token ring")
	}

	done := make(chan bool)

	tr := NewTokenRing(*chainSize, done)

	sendC := tr.Run()

	sendC <- Token{
		Data:     "Hello World",
		Reciever: 1,
		TTL:      1000,
	}

	sendC <- Token{
		Data:     "Undelivered, in some case it could be delivered :)",
		Reciever: 1000,
		TTL:      9999,
	}

	sendC <- Token{
		Data:     "Too short to live",
		TTL:      0,
		Reciever: 0,
	}

	time.Sleep(10 * time.Second)
	done <- true
}
