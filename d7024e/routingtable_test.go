package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"strconv"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	testNodes := 5
	myID := d7024e.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	rt := d7024e.NewRoutingTable(d7024e.NewNode(myID, "localhost:8000"), nil)

	for i := 1; i <= testNodes; i++ {
		address := "localhost:800" + strconv.Itoa(i+1)
		newID := d7024e.NewKademliaID("FFFFFFFF0000000000000000000000000000000" + strconv.Itoa(i))
		newNode := d7024e.NewNode(newID, address)
		rt.AddNode(newNode)
	}

	nodes := rt.FindClosestNodes(d7024e.NewKademliaID("AAAAFFFF00000000000000000000000000000000"), 20)
	count := len(nodes)
	if count != testNodes {
		t.Error(strconv.Itoa(count) + " != " + strconv.Itoa(testNodes))
	}

	//make sure the fake ID is not in the table
	fakeID := "0FFFFFFF00000000000000000000000000000000"
	fakeNodes := rt.FindClosestNodes(d7024e.NewKademliaID(fakeID), 20)
	if fakeNodes[0].ID.String() == fakeID {
		t.Error()
	}

	for index := 0; index < 160; index++ {
		id := rt.GetRandomIDInBucket(index)
		newIndex := rt.GetBucketIndex(id)
		if index != newIndex {
			t.Error(strconv.Itoa(index) + " != " + strconv.Itoa(newIndex))
		}
	}
}

func TestRemoveNodeRoutingTable(t *testing.T) {
	myNode := d7024e.NewNode(d7024e.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8000")
	otherNode := d7024e.NewNode(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	rt := d7024e.NewRoutingTable(myNode, nil)
	rt.AddNode(otherNode)
	if len(rt.FindClosestNodes(&otherNode.ID, 1)) != 1 {
		t.Error("failed to add bucket")
	}
	node := rt.FindClosestNodes(&otherNode.ID, 1)
	rt.RemoveNode(&node[0])
	if len(rt.FindClosestNodes(&otherNode.ID, 1)) != 0 {
		t.Error("failed to remove bucket")
	}
}

func TestSizeRoutingTable(t *testing.T) {
	myNode := d7024e.NewNode(d7024e.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8000")
	otherNode := d7024e.NewNode(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	rt := d7024e.NewRoutingTable(myNode, nil)
	rt.AddNode(otherNode)
	if rt.GetSize() != 1 {
		t.Error("size failed")
	}
}

func TestBucketQueue(t *testing.T) {

}
