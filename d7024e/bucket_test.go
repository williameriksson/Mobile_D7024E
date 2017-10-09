package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
)

const test_length int = 10

func TestAddNodeAndLen(t *testing.T) {
	bucket := d7024e.NewBucket()
	for i := 1; i <= test_length; i++ {
		bucket.AddNode(d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000"))
		if bucket.Len() != i {
			t.Fail()
		}
	}
}

/*
func TestGetNodeAndCalcDistance(t *testing.T) {
	bucket := d7024e.NewBucket()
	for i := 1; i <= test_length; i++ {
		bucket.AddNode(d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000"))
	}

	target := d7024e.NewRandomKademliaID()



}
*/

func TestRemoveNode(t *testing.T) {
	bucket := d7024e.NewBucket()
	node := d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000")
	bucket.AddNode(node)
	if bucket.Len() != 1 {
		t.Error("Addnode Malfuction")
	}
	bucket.RemoveNode(&node)
}

func TestQueue(t *testing.T) {
	bucket := d7024e.NewBucket()
	node := d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000")
	bucket.AddToQueue(&node)
	if bucket.QueueLen() != 1 {
		t.Error("AddToQueue Did not queue element")
	}
	nodePopped := bucket.PopQueue()
	if bucket.QueueLen() != 0 {
		t.Error("PopQueue Did not remove element")
	}
	if nodePopped != node {
		t.Error("PopQueue Did not return element")
	}
}
