package d7024e

import (
	"fmt"
	"strconv"
	"time"
)

const alpha int = 3
const k int = 20
const timeOutTime time.Duration = 1000

type Kademlia struct {
	Network       Network
	files         map[string][]byte // Hash table mapping sha-1 hash (base64 encoded) to some data
	RoutingTable  *RoutingTable
	LookupCount   int
	returnedNodes NodeCandidates
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
}

func (kademlia *Kademlia) JoinNetwork(bootStrapIP string, myIP string) {

	myID := NewRandomKademliaID()        //temp ID
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID

	myNode := NewNode(myID, myIP)
	kademlia.RoutingTable = NewRoutingTable(myNode)

	kademlia.Network = Network{&kademlia.RoutingTable.me, make(chan Message, 10), make(chan string), nil}

	conn := kademlia.Network.Listen(myIP)
	kademlia.Network.Conn = conn
	go kademlia.Network.HandleConnection()

	if bootStrapIP != "" {
		// fmt.Println("GOT bootstrap ID")
		bootStrapNode := NewNode(bootStrapID, bootStrapIP)
		go kademlia.Network.SendPingMessage(&bootStrapNode)

		//Wait for confirmation
		confirmation := <-kademlia.Network.MsgChannel

		if confirmation.Command == "PING_ACK" {
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING_ACK")
			//ping success, proceed with bootstrap procedure.
			kademlia.RoutingTable.AddNode(bootStrapNode)
			queriedNodes := make(map[string]bool)
			nodeCandidates := NodeCandidates{}
			kademlia.LookupNode(kademlia.RoutingTable.me.ID, queriedNodes, nodeCandidates, 0)
			// kademlia.Network.SendFindNodeMessage(&bootStrapNode, kademlia.RoutingTable.me.ID)

		} else {
			fmt.Println("Failed connect!")
		}
	}
	kademlia.channelReader()

}

func (kademlia *Kademlia) channelReader() {
	for {
		// kademlia.Network.TestChannel <- "chan read"
		// time.Sleep(time.Millisecond * 100)
		msg := <-kademlia.Network.MsgChannel

		switch msg.Command {
		case "PING_ACK":
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING_ACK")
		case "PING":
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING")
			kademlia.Network.SendPingAck(msg.SenderNode)
		case "STORE":
			fmt.Println("GOT STORE")
			kademlia.Store(msg.Data)
		case "FIND_NODE":
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT FIND_NODE")
			// A FIND_NODE msg recived, send back k closest nodes.
			kademliaID := NewKademliaID(string(msg.Data))
			nodeList := kademlia.RoutingTable.FindClosestNodes(kademliaID, 20)
			kademlia.RoutingTable.AddNode(*msg.SenderNode)
			kademlia.Network.SendReturnFindNodeMessage(msg.SenderNode, nodeList)

			temp := kademlia.RoutingTable.me.Address + "\n"
			for j := 1; j <= len(kademlia.RoutingTable.buckets); j++ {
				if kademlia.RoutingTable.buckets[j-1].Len() > 0 {
					temp += strconv.Itoa(kademlia.RoutingTable.buckets[j-1].Len())
				}
				temp += " "
				if j%20 == 0 {
					temp += "\n"
				}
			}
			kademlia.Network.TestChannel <- temp

		case "RETURN_FIND_NODE":
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT RETURN_FIND_NODE")
			kademlia.returnedNodes.Append(msg.NodeList)
			kademlia.LookupCount++
			for i := 0; i < len(msg.NodeList); i++ {
				kademlia.RoutingTable.AddNode(msg.NodeList[i])
			}

			temp := kademlia.RoutingTable.me.Address + "\n"
			for j := 1; j <= len(kademlia.RoutingTable.buckets); j++ {
				if kademlia.RoutingTable.buckets[j-1].Len() > 0 {
					temp += strconv.Itoa(kademlia.RoutingTable.buckets[j-1].Len())
				}
				temp += " "
				if j%20 == 0 {
					temp += "\n"
				}
			}
			kademlia.Network.TestChannel <- temp

		case "FIND_VALUE":
			fmt.Println("FIND_VALUE")
			kademlia.LookupValue(msg.Hash)
		default:
			fmt.Println("GOT DEFAULT")
		}
	}
}

func (kademlia *Kademlia) LookupNode(target *KademliaID, queriedNodes map[string]bool, prevBestNodes NodeCandidates, recCount int) {
	// Should be run when looking for a node (not during bootstrap though)
	kademlia.LookupCount = 0
	kademlia.returnedNodes = prevBestNodes

	fmt.Println("LookupNode running")
	closestNodes := kademlia.RoutingTable.FindClosestNodes(target, k)
	for i := 0; i < alpha && i < len(closestNodes); i++ {
		if queriedNodes[target.String()] == false {
			queriedNodes[target.String()] = true
			kademlia.Network.SendFindNodeMessage(&closestNodes[i], target)
		} else {

		}
	}

	timeout := false
	timeStamp := time.Now()

	for kademlia.LookupCount < alpha && (kademlia.LookupCount < len(closestNodes) && len(closestNodes) < alpha) && !timeout {
		//busy waiting for RETURN_FIND_NODE
		if time.Now().Sub(timeStamp) > timeOutTime {
			timeout = true
		}
	}

	if !timeout {
		kademlia.returnedNodes.Sort() //what does this do
		bestNodes := NodeCandidates{nodes: kademlia.returnedNodes.GetNodes(20)}
		if bestNodes.nodes[0].ID.String() == target.String() {
			//first node IS target means we DID find it this run.
		} else if recCount == k {
			//did NOT find node after k attempts

		} else {
			//did NOT find node, continue search
			kademlia.LookupNode(target, queriedNodes, bestNodes, 0)
		}
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
