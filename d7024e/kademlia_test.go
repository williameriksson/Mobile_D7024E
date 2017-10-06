package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"fmt"
	"strconv"
	"testing"
	"time"
)

const timeOut time.Duration = (2 * time.Second)

func TestKademliaInstantiation(t *testing.T) {
	kad1 := d7024e.NewKademlia()
	// kad2 := d7024e.NewKademlia()
	// kad3 := d7024e.NewKademlia()
	// kad4 := d7024e.NewKademlia()

	if kad1.LookupCount != 0 || kad1.RoutingTable != nil {
		t.Fail()
	}
}

// func TestStore(t *testing.T) {
// 	// THIS IS NOT FINAL, SHOULD BE MADE TO BLOCK INSTEAD OF SLEEP
// 	kademlia1 := d7024e.NewKademlia()
// 	go kademlia1.Run("", "127.0.0.1:8000")
//
// 	kademlia2 := d7024e.NewKademlia()
// 	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8002")
//
// 	time.Sleep(time.Millisecond * 500)
//
// 	data := []byte("Test String")
// 	hash := d7024e.HashData(data)
// 	nodes := kademlia2.RoutingTable.FindClosestNodes(d7024e.NewKademliaID(hash), 1)
// 	//fmt.Printf("node: %x\n", nodes[0].ID)
// 	kademlia2.Network.SendStoreMessage(&nodes[0], data)
//
// 	time.Sleep(time.Millisecond * 500)
//
// 	kademlia1.LookupValue(hash)
//
// 	time.Sleep(time.Millisecond * 500)
// }

/*
*	Sets up a network of n nodes, conditons for PASS is that all nodes in the network has knowledge of at least k(20) nodes
 */
func TestKademlia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping due to -short flag")
	}
	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 100)

	var port int = 8001
	count := 100

	var kademliaArray [100]*d7024e.Kademlia
	progress := 0

	for i := 0; i < count; i++ {
		kademliaArray[i] = d7024e.NewKademlia()
		go kademliaArray[i].Run("127.0.0.1:8000", "127.0.0.1:"+strconv.Itoa(port+i))
		time.Sleep(time.Millisecond * 100)
		print("\033[H\033[2J")
		progress++
		for c := 0; c < 100; c++ {
			if c <= progress {
				fmt.Print("|")
			} else {
				fmt.Print("-")
			}
		}
	}
	fmt.Println("\nProcessing...")
	var minBucketPop = count
	if minBucketPop > 20 {
		minBucketPop = 20
	}

	time.Sleep(time.Millisecond * 2000)
	if len(kademlia1.RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)) < minBucketPop {
		t.Error("Bootstrap fail")
	}

	for j := 0; j < count; j++ {
		if len(kademliaArray[j].RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)) < minBucketPop {
			t.Error(kademliaArray[j].RoutingTable.GetMyAdress() + " : " + strconv.Itoa(len(kademliaArray[j].RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), count))) + " != " + strconv.Itoa(minBucketPop))
		} else {
			// fmt.Println(kademliaArray[j].RoutingTable.GetMyAdress() + " : " + strconv.Itoa(len(kademliaArray[j].RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), count))))
		}
	}
}

func TestPingTimeout(t *testing.T) {
	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 100)

	kademlia2 := d7024e.NewKademlia()
	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8001")
	time.Sleep(time.Millisecond * 100)

	kad1rtNodes := kademlia1.RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)
	// kad2rtNodes := kademlia2.RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)

	kademlia1.RoutingTable.CheckAlive(kad1rtNodes)

	kad1rtNodes = kademlia1.RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)
	// kad2rtNodes = kademlia2.RoutingTable.FindClosestNodes(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), 20)

	// time.Sleep(time.Millisecond * 3000)

	//no real way of testing this.
	// if len(kad1rtNodes) != 1 {
	// 	t.Error("fail")
	// }
	// if len(kad2rtNodes) != 1 {
	// 	t.Error("fail")
	// }
}
