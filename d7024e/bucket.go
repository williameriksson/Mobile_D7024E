package d7024e

import (
	"container/list"
	"fmt"
)

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

func (bucket *bucket) AddNode(node Node) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Node).ID

		if (node).ID.Equals(&nodeID) {
			element = e
			break
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(node)
			fmt.Printf("PUSHING FRONT: 0x%X\n", node.ID)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) GetNodeAndCalcDistance(target *KademliaID) []Node {
	var nodes []Node

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		node := elt.Value.(Node)
		node.CalcDistance(target)
		nodes = append(nodes, node)
	}

	return nodes
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
