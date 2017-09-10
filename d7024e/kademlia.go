package d7024e

import "fmt"

const alpha int = 3

type Kademlia struct {
	network Network
	files map[string][]byte  // Hash table mapping sha-1 hash (base64 encoded) to some data
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
		fmt.Println(<-kademlia.network.testChannel)
	}
}

func (kademlia *Kademlia) JoinNetwork(bootStrapIP string, myIP string) {

	myID := NewRandomKademliaID()        //temp ID
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID

	myNode := NewNode(myID, myIP)
	routingTable := NewRoutingTable(myNode)

	kademlia.network = Network{&routingTable.me, make(chan Message), make(chan string)}

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
		kademlia.network.testChannel <- "chan read"
		<-kademlia.network.msgChannel
	}
}

func (kademlia *Kademlia) LookupNode(target *Node) {
	// TODO
	fmt.Println("LookupNode running")
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
	fmt.Println(kademlia.files[hash])
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
	hash := HashData(data)
	kademlia.files[hash] = data
}