package d7024e

import "fmt"

const alpha int = 3

type Kademlia struct {
	Network Network
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
	routingTable := NewRoutingTable(myNode)

	kademlia.Network = Network{&routingTable.me, make(chan Message), make(chan string), nil}

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
			routingTable.AddNode(bootStrapNode)
			kademlia.LookupNode(&routingTable.me)
		} else {
			fmt.Println("Failed connect!")
		}
	}
	go kademlia.channelReader()

}

func (kademlia *Kademlia) channelReader() {
	for {
		kademlia.Network.TestChannel <- "chan read"
		<-kademlia.Network.MsgChannel
	}
}

func (kademlia *Kademlia) LookupNode(target *Node) {
	// TODO
	fmt.Println("LookupNode running")
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
