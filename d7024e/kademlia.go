package d7024e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"strconv"
	"time"
)

const alpha int = 3
const k int = 20
const timeOutTime time.Duration = 1000

type Kademlia struct {
	Network       Network
	files         map[string]string // Hash table mapping sha-1 hash (base64 encoded) to some data
	RoutingTable  *RoutingTable
	LookupCount   int
	LookupValueCount int
	returnedNodes NodeCandidates
	returnedValueNodes NodeCandidates
	returnedValue []byte
	Datainfo	DataInformation
}

// Constructor
func NewKademlia() *Kademlia {
	var kademlia Kademlia
	kademlia.files = make(map[string]string)
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

	myID := NewRandomKademliaID() //temp ID
	//fmt.Printf("ID: 0x%X\n", myID)
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID TODO: bootstrap should NOT be assigned a random ID.

	myNode := NewNode(myID, myIP)
	kademlia.RoutingTable = NewRoutingTable(myNode)

	kademlia.Network = Network{me: &kademlia.RoutingTable.me, MsgChannel: make(chan Message), TestChannel: make(chan string, 100)}
	// kademlia.Network.TestChannel <- ("My ID : " + myID.String())


	go kademlia.RepublishData()
	go kademlia.PurgeData()
	conn := kademlia.Network.Listen(myIP)
	kademlia.Network.Conn = conn
	go kademlia.Network.HandleConnection()

	/*	--- Bootstrap Procedure ---	*/
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

				secondConfirm := <-kademlia.Network.MsgChannel
				if secondConfirm.Command == cmd_find_node_returned {
					var nodeList []Node
					err := json.Unmarshal(secondConfirm.Data, &nodeList)
					checkError(err)
					kademlia.findNodeReturn(&secondConfirm.SenderNode, nodeList)
					kademlia.RefreshBuckets()
					kademlia.PingAllNodes()
					break
				}

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
const cmd_value_returned = "VALUE_RETURNED"

// TODO: FIX THE PATH OS INDEPENDETN
const DOWNLOAD_PATH = "./Downloads"

func (kademlia *Kademlia) channelReader() {
	for {
		//halts here waiting for a command.
		msg := <-kademlia.Network.MsgChannel
		kademlia.RoutingTable.AddNode(msg.SenderNode)
		switch msg.Command {
		case cmd_ping_ack:
			kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT PING_ACK")

		case cmd_ping:
			kademlia.Network.SendPingAck(&msg.SenderNode)

		case cmd_store:
			fmt.Println("GOT " + cmd_store)
			// TODO: Add call to own server to establish tcp conn and get the actual file
			kademlia.Store(string(msg.Data), DOWNLOAD_PATH, false)

		case cmd_find_node:
			//THIS node has recived a request to find a certain node (from some other node)
			// kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT FIND_NODE")
			var kID KademliaID
			err := json.Unmarshal(msg.Data, &kID)
			checkError(err)
			kademlia.findNode(&msg.SenderNode, &kID)

		case cmd_find_node_returned:
			//Some node has returned a list of the k closest nodes to the node that THIS node requested
			// kademlia.Network.TestChannel <- kademlia.RoutingTable.me.Address + (" GOT RETURN_FIND_NODE")
			var nodeList []Node
			err := json.Unmarshal(msg.Data, &nodeList)
			checkError(err)
			kademlia.findNodeReturn(&msg.SenderNode, nodeList)

		case cmd_find_value:
			//THIS node has recived a request to find a certain value (from some other node)
			fmt.Println("FIND_VALUE")
			var kID KademliaID
			err := json.Unmarshal(msg.Data, &kID)
			checkError(err)
			kademlia.FindValue(&msg.SenderNode, &kID)

		case cmd_find_value_returned:
			//Some node has returned a list of the k closest nodes to the value that THIS node requested
			var nodeList []Node
			err := json.Unmarshal(msg.Data, &nodeList)
			checkError(err)
			kademlia.FindValueReturn(&msg.SenderNode, nodeList)

		case cmd_value_returned:
			fmt.Println("Got a value returned to me!")
			//Some node has returned the value that THIS node requested
			kademlia.returnedValue = msg.Data
			kademlia.Store(msg.Data, false)

		default:
			fmt.Println("GOT DEFAULT")
		}
		tmpString := ""

		for key, value := range kademlia.files {
	    tmpString += "Key: " + key + "Value: " + string(value)
	  }


		kademlia.WriteToFile(kademlia.RoutingTable.me.ID.String()+".txt", []byte(tmpString))
		//kademlia.WriteToFile(kademlia.RoutingTable.me.ID.String()+".txt", []byte(kademlia.RoutingTable.GetRoutingTable()))
	}
}

func (kademlia *Kademlia) findNode(senderNode *Node, kID *KademliaID) {
	nodeList := kademlia.RoutingTable.FindClosestNodes(kID, k)
	//fmt.Println("nodelist --")
	// for i := 0; i < len(nodeList); i++ {
	// 	fmt.Printf("node : 0x%X\n", nodeList[i].ID)
	// }
	kademlia.Network.SendReturnFindNodeMessage(senderNode, nodeList)
}

func (kademlia *Kademlia) findNodeReturn(senderNode *Node, nodeList []Node) {
	kademlia.returnedNodes.Append(nodeList)
	kademlia.LookupCount++

	//adds all the returned nodes to the RoutingTable
	for i := 0; i < len(nodeList); i++ {
		// kademlia.Network.SendPingMessage(&nodeList[i])
		kademlia.RoutingTable.AddNode(nodeList[i])
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
		}
	}

	timeout := false
	timeStamp := time.Now()

	for (kademlia.LookupCount < 1) && !timeout {
		//busy waiting for at least one RETURN_FIND_NODE
		if time.Now().Sub(timeStamp) > timeOutTime {
			timeout = true
		}
	}

	if !timeout {
		fmt.Println("RETURNEDNODES: ", kademlia.returnedValueNodes)
		kademlia.returnedNodes.Sort() //what does this do
		bestNodes := NodeCandidates{nodes: kademlia.returnedNodes.GetNodes(k)}
		if bestNodes.nodes[0].ID.String() == target.String() {
			//first node IS target means we DID find it this run.
		} else if recCount == k {
			//did NOT find node after k attempts

		} else {
			//did NOT find node, continue search
			kademlia.LookupNode(target, queriedNodes, bestNodes, (recCount + 1))
		}
	}
}

func (kademlia *Kademlia) RefreshBuckets() {
	myIndex := kademlia.RoutingTable.GetBucketIndex(&kademlia.RoutingTable.me.ID)
	//for the buckets less than "me"
	for i := myIndex; i >= 0; i-- {
		if kademlia.RoutingTable.buckets[i].Len() < 1 {
			kadID := kademlia.RoutingTable.GetRandomIDInBucket(i)
			receiverNode := kademlia.RoutingTable.FindClosestNodes(kadID, 1)
			kademlia.Network.SendFindNodeMessage(&receiverNode[0], kadID)
		}
	}
	//for the buckets more than "me"
	for j := myIndex; j < (IDLength * 8); j++ {
		if kademlia.RoutingTable.buckets[j].Len() < 1 {
			kadID := kademlia.RoutingTable.GetRandomIDInBucket(j)
			receiverNode := kademlia.RoutingTable.FindClosestNodes(kadID, 1)
			kademlia.Network.SendFindNodeMessage(&receiverNode[0], kadID)
		}
	}
}

func (kademlia *Kademlia) PingAllNodes() {
	nodes := kademlia.RoutingTable.FindClosestNodes(&kademlia.RoutingTable.me.ID, kademlia.RoutingTable.GetSize())
	for i := 0; i < len(nodes); i++ {
		kademlia.Network.SendPingMessage(&nodes[i])
	}
}

func (kademlia *Kademlia) PublishData(hash string, path string) {

	closestNodes := kademlia.RoutingTable.FindClosestNodes(NewKademliaID(hash), k)
	for i := 0; i < len(closestNodes); i++ {
    kademlia.Network.SendStoreMessage(&closestNodes[i], []byte(hash))
	}
	kademlia.Store(hash, path, true)
}

func (kademlia *Kademlia) Get(hash string) string{
	return kademlia.files[hash]
}

/*func (kademlia *Kademlia) Store(hash string, path string){
	kademlia.files[hash] = path
}*/

func (kademlia *Kademlia) Store(hash string, path string, me bool) {
	//hash := HashStr(fileName)
	kademlia.files[hash] = path

	if me {
		for _, myKey := range kademlia.Datainfo.MyKeys {
			if (myKey == hash) {
				return
			}
		}
		kademlia.Datainfo.MyKeys = append(kademlia.Datainfo.MyKeys, hash)
		return
	}

	for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
		if (purgeInfo.Key == hash) {
			kademlia.SetPurgeStamp(&purgeInfo)
			return
		}
	}

	newPurgeInfo := PurgeInformation{Key: hash, Pinned: false}
	kademlia.SetPurgeStamp(&newPurgeInfo)
	kademlia.Datainfo.PurgeInfos = append(kademlia.Datainfo.PurgeInfos, newPurgeInfo)

}

func (kademlia *Kademlia) WriteToFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("WRITE ERROR: ", err)
	}
}
