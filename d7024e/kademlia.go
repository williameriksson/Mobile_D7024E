package d7024e

import "fmt"

const alpha int = 3

type Kademlia struct {
	network Network
}

/*
* Initializes the bootstrap procedure
 */
func (kademlia *Kademlia) Run(connectIP string, myIP string) {
	bootStrapIPtemp := connectIP
	myIPtemp := myIP

	kademlia.JoinNetwork(bootStrapIPtemp, myIPtemp) //attempts joining the Kademlia network
}

func (kademlia *Kademlia) JoinNetwork(bootStrapIP string, myIP string) {

	myID := NewRandomKademliaID()        //temp ID
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID

	myNode := NewNode(myID, myIP)
	routingTable := NewRoutingTable(myNode)

	kademlia.network = Network{&routingTable.me, make(chan Message)}

	go kademlia.network.Listen(myIP)

	if bootStrapIP != "" {
		fmt.Println("GOT bootstrap ID")
		bootStrapNode := NewNode(bootStrapID, bootStrapIP)
		go kademlia.network.SendPingMessage(&bootStrapNode)

		//Wait for confirmation
		confirmation := <-kademlia.network.msgChannel

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
		<-kademlia.network.msgChannel
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
