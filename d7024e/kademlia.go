package d7024e

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"io/ioutil"
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
	fmt.Println("ID: ", myID)
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID TODO: bootstrap should NOT be assigned a random ID.

	myNode := NewNode(myID, myIP)
	kademlia.RoutingTable = NewRoutingTable(myNode)

	kademlia.Network = Network{me: &kademlia.RoutingTable.me, MsgChannel: make(chan Message), TestChannel: make(chan string, 100)}
	kademlia.Network.TestChannel <- ("My ID : " + myID.String())

	conn := kademlia.Network.Listen(myIP)
	kademlia.Network.Conn = conn
	go kademlia.Network.HandleConnection()

	if bootStrapIP != "" {
		// fmt.Println("GOT bootstrap ID")
		bootStrapNode := NewNode(bootStrapID, bootStrapIP)
		kademlia.Network.SendPingMessage(&bootStrapNode)

		//Wait for confirmation
		for {
			confirmation := <-kademlia.Network.MsgChannel

			if confirmation.Command == cmd_ping_ack {
				kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING_ACK")
				//ping success, proceed with bootstrap procedure.
				//kademlia.RoutingTable.AddNode(*confirmation.SenderNode) // QUESTION: Why does this line produce duplicates but the line below does not?
				kademlia.RoutingTable.AddNode(NewNode(NewKademliaID(confirmation.SenderNode.ID.String()), bootStrapIP))
				queriedNodes := make(map[string]bool)
				nodeCandidates := NodeCandidates{}
				kademlia.LookupNode(&kademlia.RoutingTable.me.ID, queriedNodes, nodeCandidates, 0)
				// kademlia.Network.SendFindNodeMessage(&bootStrapNode, kademlia.RoutingTable.me.ID) // QUESTION: Why was this commented out?
				break
				} else {
					fmt.Println("Expected PING_ACK, instead got: ", confirmation.Command, "will keep waiting for PING_ACK")
				}

		}
	}
	kademlia.channelReader()

}

//the KADEMLIA commands
const cmd_ping = "PING"
const cmd_store = "STORE"
const cmd_find_node = "FIND_NODE"
const cmd_find_value = "FIND_VALUE"

//response commands
const cmd_ping_ack = "PING_ACK"
const cmd_find_node_returned = "FIND_NODE_RETURNED"
const cmd_find_value_returned = "FIND_VALUE_RETURNED"

func (kademlia *Kademlia) channelReader() {
	for {
		//halts here waiting for a command.
		msg := <-kademlia.Network.MsgChannel
		kademlia.RoutingTable.AddNode(msg.SenderNode) // QUESTION: Added this line, should this be here?
		switch msg.Command {
		case cmd_ping_ack:
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING_ACK")

		case cmd_ping:
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING")
			kademlia.Network.SendPingAck(&msg.SenderNode)

		case cmd_store:
			fmt.Println("GOT " + cmd_store)
			kademlia.Store(msg.Data)

		case cmd_find_node:
			//THIS node has recived a request to find a certain node (from some other node)
			// kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT FIND_NODE")
			kID := NewKademliaID(string(msg.Data))
			kademlia.findNode(&msg.SenderNode, kID)

		case cmd_find_node_returned:
			//Some node has returned a list of the k closest nodes to the node that THIS node requested
			// kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT RETURN_FIND_NODE")
			var nodeList []Node
			err := json.Unmarshal(msg.Data, &nodeList)
			checkError(err)
			kademlia.findNodeReturn(&msg.SenderNode, nodeList)

		case cmd_find_value:
			fmt.Println("FIND_VALUE")
			kademlia.LookupValue(msg.Hash)

		default:
			fmt.Println("GOT DEFAULT")
		}
		kademlia.WriteToFile(kademlia.RoutingTable.me.ID.String(), []byte(kademlia.RoutingTable.GetRoutingTable()))
	}
}

func (kademlia *Kademlia) findNode(senderNode *Node, kID *KademliaID) {
	kademlia.RoutingTable.AddNode(*senderNode)
	nodeList := kademlia.RoutingTable.FindClosestNodes(kID, k)
	fmt.Println("nodelist --")
	for i := 0; i < len(nodeList); i++ {
		fmt.Println("node : " + nodeList[i].ID.String())
	}
	kademlia.Network.SendReturnFindNodeMessage(senderNode, nodeList)

	// kademlia.Network.TestChannel <- ("sending back : " + strconv.Itoa(len(nodeList)))
	// temp := kademlia.RoutingTable.me.Address + "\n"
	// for j := 1; j <= len(kademlia.RoutingTable.buckets); j++ {
	// 	if kademlia.RoutingTable.buckets[j-1].Len() > 0 {
	// 		temp += strconv.Itoa(kademlia.RoutingTable.buckets[j-1].Len())
	// 	}
	// 	temp += " "
	// 	if j%20 == 0 {
	// 		temp += "\n"
	// 	}
	// }
	// kademlia.Network.TestChannel <- temp
}

func (kademlia *Kademlia) findNodeReturn(senderNode *Node, nodeList []Node) {
	kademlia.returnedNodes.Append(nodeList)
	kademlia.LookupCount++

	//adds all the returned nodes to the RoutingTable
	for i := 0; i < len(nodeList); i++ {
		kademlia.RoutingTable.AddNode(nodeList[i])
	}
	kademlia.Network.TestChannel <- ("recieved nodes : " + strconv.Itoa(len(nodeList)))

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
}

func (kademlia *Kademlia) LookupNode(target *KademliaID, queriedNodes map[string]bool, prevBestNodes NodeCandidates, recCount int) {
	// Should be run when looking for a node (not during bootstrap though)
	kademlia.LookupCount = 0
	kademlia.returnedNodes = prevBestNodes

	fmt.Println("LookupNode running")
	closestNodes := kademlia.RoutingTable.FindClosestNodes(target, k)
	fmt.Println("closest len " + strconv.Itoa(len(closestNodes)))
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
		bestNodes := NodeCandidates{nodes: kademlia.returnedNodes.GetNodes(k)}
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

/*
func (kademlia *Kademlia) LookupNode2(target *KademliaID) {
	closestNodes []Node = kademlia.RoutingTable.FindClosestNodes(target, alpha)
	for _, node := range closestNodes {
		kademlia.Network.SendFindNodeMessage(node, target)
	}

}*/

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
		id := NewKademliaID(hash)
		//kademlia.LookupNode
		fmt.Println("No value, %x", id)
	}

}

func (kademlia *Kademlia) Store(data []byte) {
	hash := HashData(data)
	kademlia.files[hash] = data
}

func (kademlia *Kademlia) WriteToFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("WRITE ERROR: ", err)
	}
}
