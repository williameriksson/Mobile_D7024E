package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
<<<<<<< HEAD
	//"fmt"
=======
>>>>>>> refs/remotes/origin/master
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

func TestStore(t *testing.T) {
	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")

	time.Sleep(time.Millisecond * 1000)

	kademlia2 := d7024e.NewKademlia()
	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8002")

	time.Sleep(time.Millisecond * 1000)

	data := []byte("Test String")
	hash := d7024e.HashData(data)
	nodes := kademlia2.RoutingTable.FindClosestNodes(d7024e.NewKademliaID(hash), 1)
	//fmt.Printf("node: %x\n", nodes[0].ID)
	kademlia2.Network.SendStoreMessage(&nodes[0], data)

	time.Sleep(time.Millisecond * 1000)

	kademlia1.LookupValue(hash)

	time.Sleep(time.Millisecond * 1000)
}
/*
func TestKademlia(t *testing.T) {
	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")

	time.Sleep(time.Millisecond * 1000)

	kademlia2 := d7024e.NewKademlia()
	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8002")

	time.Sleep(time.Millisecond * 1000)

	kademlia3 := d7024e.NewKademlia()
	go kademlia3.Run("127.0.0.1:8000", "127.0.0.1:8003")

	time.Sleep(time.Millisecond * 1000)

	kademlia4 := d7024e.NewKademlia()
	go kademlia4.Run("127.0.0.1:8000", "127.0.0.1:8004")

	time.Sleep(time.Millisecond * 1000)

	kademlia5 := d7024e.NewKademlia()
	go kademlia5.Run("127.0.0.1:8000", "127.0.0.1:8005")

	time.Sleep(time.Millisecond * 1000)
	
	for {
		select {
		case msg := <-kademlia1.Network.TestChannel:
			fmt.Println(msg)
		case msg := <-kademlia2.Network.TestChannel:
			fmt.Println(msg)
		case msg := <-kademlia3.Network.TestChannel:
			fmt.Println(msg)
		case msg := <-kademlia4.Network.TestChannel:
			fmt.Println(msg)
		case msg := <-kademlia5.Network.TestChannel:
			fmt.Println(msg)
		}
	}
}*/