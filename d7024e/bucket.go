package d7024e

import (
	"container/list"
)

type bucket struct {
	list  *list.List
	queue *list.List
}

func NewBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.queue = list.New()
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
			//fmt.Printf("PUSHING FRONT: 0x%X\n", node.ID)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) RemoveNode(node *Node) {
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Node).ID

		if (node).ID.Equals(&nodeID) {
			bucket.list.Remove(e)
		}
	}
}

//adds a node to the queue (the list of nodes to be added when k-list isn't full anymore)
func (bucket *bucket) AddToQueue(node *Node) {
	bucket.queue.PushBack(*node)
}

func (bucket *bucket) PopQueue() Node {
	var node Node
	node = bucket.queue.Front().Value.(Node)
	bucket.queue.Remove(bucket.queue.Front())
	return node //throw an error perhaps?
}

func (bucket *bucket) GetNodelist() []Node {
	var nodes []Node

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		nodes = append(nodes, node)
	}
	return nodes
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

func (bucket *bucket) QueueLen() int {
	return bucket.queue.Len()
}
