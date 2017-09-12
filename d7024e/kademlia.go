package d7024e

import (
	"fmt"
)

const alpha int = 3

type Kademlia struct {
	Network      Network
	files        map[string][]byte // Hash table mapping sha-1 hash (base64 encoded) to some data
	RoutingTable *RoutingTable
}

// Constructor
func NewKademlia() *Kademlia {
	var kademlia Kademlia
	kademlia.files = make(map[string][]byte)
	return &kademlia
}

/*
* Initializes the bootstrap procedure
 */

func (kademlia *Kademlia) Run(connectIP string, myIP string) {
	bootStrapIPtemp := connectIP
	myIPtemp := myIP

	kademlia.JoinNetwork(bootStrapIPtemp, myIPtemp) //attempts joining the Kademlia network
	for {
	}
}

func (kademlia *Kademlia) JoinNetwork(bootStrapIP string, myIP string) {

	myID := NewRandomKademliaID()        //temp ID
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID

	myNode := NewNode(myID, myIP)
	kademlia.RoutingTable = NewRoutingTable(myNode)

	kademlia.Network = Network{&kademlia.RoutingTable.me, make(chan Message), make(chan string), nil}

	conn := kademlia.Network.Listen(myIP)
	kademlia.Network.Conn = conn
	go kademlia.Network.HandleConnection()

	if bootStrapIP != "" {
		fmt.Println("GOT bootstrap ID")
		bootStrapNode := NewNode(bootStrapID, bootStrapIP)
		go kademlia.Network.SendPingMessage(&bootStrapNode)

		//Wait for confirmation
		confirmation := <-kademlia.Network.MsgChannel
		fmt.Println(confirmation)

		if confirmation.Command == "PING_ACK" {
			//ping success, proceed with bootstrap procedure.
			kademlia.RoutingTable.AddNode(bootStrapNode)
			kademlia.LookupNode(&kademlia.RoutingTable.me)
		} else {
			fmt.Println("Failed connect!")
		}
	}
	go kademlia.channelReader()

}

func (kademlia *Kademlia) channelReader() {
	for {
		kademlia.Network.TestChannel <- "chan read"
		msg := <-kademlia.Network.MsgChannel

		switch msg.Command {
		case "PING_ACK":
			fmt.Println("GOT PING_ACK")
		case "PING":
			fmt.Println(kademlia.RoutingTable.me.ID.String() + "GOT PING")
			kademlia.Network.SendPingAck(msg.SenderNode)
		case "STORE":
			fmt.Println("GOT STORE")
		case "FIND_NODE":
			fmt.Println("GOT FIND_NODE")
		case "FIND_VALUE":
			fmt.Println("FIND_VALUE")
		default:
			fmt.Println("GOT DEFAULT")
		}
	}
}

func (kademlia *Kademlia) LookupNode(target *Node) {
	// TODO
	fmt.Println("LookupNode running")
	closestNodes := kademlia.RoutingTable.FindClosestNodes(target.ID, alpha)
	for i := 0; i < alpha && i < len(closestNodes); i++ {
		msg := Message{"FIND_NODE", kademlia.Network.me}
		go kademlia.Network.sendMessage(&closestNodes[i], &msg)
	}
}

func (kademlia *Kademlia) LookupValue(hash string) {
	// If the node has the value, return it
	if val, ok := kademlia.files[hash]; ok {
		/*
		 * DATA (val) SHOULD BE RETURNED HERE
		 */
		fmt.Printf("Yes, the value is %x \n", val)
	} else {
		/* 
		 * TRIPLE (IP Adress, UDP Port, Node ID) SHOULD BE RETURNED HERE
		 */
		fmt.Println("No value")
	}
	
}

func (kademlia *Kademlia) Store(data []byte) {
	hash := HashData(data)
	kademlia.files[hash] = data
}
