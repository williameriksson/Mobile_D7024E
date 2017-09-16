package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := d7024e.NewRoutingTable(d7024e.NewNode(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	rt.AddNode(d7024e.NewNode(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	nodes := rt.FindClosestNodes(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), 20)
	count := len(nodes)
	if count != 6 {
		t.Fail()
	}

	//make sure the fake ID is not in the table
	fakeID := "2111111400000000000000000000000000000123"
	fakeNodes := rt.FindClosestNodes(d7024e.NewKademliaID(fakeID), 20)
	if fakeNodes[0].ID.String() == fakeID {
		t.Fail()
	}
}
