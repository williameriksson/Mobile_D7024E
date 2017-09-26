package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"strconv"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	testNodes := 5

	rt := d7024e.NewRoutingTable(d7024e.NewNode(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

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
}
