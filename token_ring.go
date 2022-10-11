// Copyright 2022 Lipatov Alexander

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "log"

// Token is a message type.
type Token struct {
	Data     string `json:"data"`
	Reciever int    `json:"reciever"`
	TTL      int    `json:"ttl"`
}

type TokenRing struct {
	Nodes []*Node
}

func NewTokenRing(size int) *TokenRing {
	tr := &TokenRing{
		Nodes: make([]*Node, 0, size),
	}

	if size < 2 {
		log.Fatal("Expected to have at least 2 size of ring")
	}

	firstNode := &Node{
		ID:    0,
		NextC: make(chan Token),
	}
	tr.Nodes = append(tr.Nodes, firstNode)

	for i := 1; i < size; i++ {
		tr.Nodes = append(tr.Nodes, &Node{
			ID:      i,
			BeforeC: tr.Nodes[i-1].NextC,
			NextC:   make(chan Token),
		})
	}
	// Assign firstNode beforeC from last node
	firstNode.BeforeC = tr.Nodes[size-1].NextC

	return tr
}

func (tr *TokenRing) Run() chan Token {
	for _, node := range tr.Nodes {
		go node.Run()
	}
	return tr.Nodes[len(tr.Nodes)/2].NextC
}

// Node represents a single block in a ring.
type Node struct {
	ID int

	BeforeC <-chan Token
	NextC   chan Token
}

// Run runs blocking listening for token.
func (node *Node) Run() {
	for t := range node.BeforeC {
		node.process(t)
	}
}

func (node *Node) process(t Token) {
	switch {
	case t.Reciever == node.ID:
		log.Printf("Token has been accepted by %d; message: %s (with left ttl = %d)", t.Reciever, t.Data, t.TTL)
	case t.TTL > 0:
		t.TTL -= 1
		node.NextC <- t
	default:
		log.Printf("Token for %d is expired", t.Reciever)
	}
}
