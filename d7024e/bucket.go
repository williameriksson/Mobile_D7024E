package d7024e

import (
	"container/list"
	"sync"
)

type bucket struct {
	list  *list.List
	queue *list.List
	mutex	sync.Mutex
}

func NewBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.queue = list.New()
	return bucket
}

func (bucket *bucket) AddNode(node Node) bool {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
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
		} else {
			return false //the bucket is full and new node cannot be inserted
		}
	} else {
		bucket.list.MoveToFront(element)
	}
	return true //either new node inserted or old node moved forward
}

func (bucket *bucket) RemoveNode(node *Node) {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Node).ID

		if (node).ID.Equals(&nodeID) {
			bucket.list.Remove(e)
		}
	}
}

//adds a node to the queue (the list of nodes to be added when k-list isn't full anymore)
func (bucket *bucket) AddToQueue(node *Node) {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.queue.PushBack(*node)
}

func (bucket *bucket) PopQueue() Node {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	var node Node
	node = bucket.queue.Front().Value.(Node)
	bucket.queue.Remove(bucket.queue.Front())
	return node //throw an error perhaps?
}

func (bucket *bucket) GetNodelist() []Node {
	var nodes []Node
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		nodes = append(nodes, node)
	}
	return nodes
}

func (bucket *bucket) GetNodeAndCalcDistance(target *KademliaID) []Node {
	var nodes []Node
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

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
