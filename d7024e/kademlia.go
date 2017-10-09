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
const timeOutTime time.Duration = time.Duration(40) * time.Millisecond

type Kademlia struct {
	Network            Network
	files              map[string]string // Hash table mapping sha-1 hash (base64 encoded) to some data
	RoutingTable       *RoutingTable
	LookupCount        int
	LookupValueCount   int
	returnedNodes      NodeCandidates
	returnedValueNodes NodeCandidates
	foundHashes        map[string]Node
	Datainfo           DataInformation
	pingedNodes        map[Node]bool
	timeoutChannel     chan bool
	valueTimeoutChan   chan bool
	ServerChannel      chan Handle
}

// Constructor
func NewKademlia() *Kademlia {
	var kademlia Kademlia
	kademlia.files = make(map[string]string)
	kademlia.pingedNodes = make(map[Node]bool)
	kademlia.timeoutChannel = make(chan bool)
	kademlia.valueTimeoutChan = make(chan bool)
	kademlia.ServerChannel = make(chan Handle)
	kademlia.foundHashes = make(map[string]Node)
	kademlia.Datainfo.PurgeInfos = make(map[string]PurgeInformation)
	kademlia.Datainfo.MyKeys = make(map[string]bool)
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
	bootStrapID := NewRandomKademliaID() //20 byte id temp ID (to allow using bootstrap as a node, discarded later)
	myNode := NewNode(myID, myIP)

	kademlia.Network = Network{MsgChannel: make(chan Message), TestChannel: make(chan string, 100)}
	kademlia.RoutingTable = NewRoutingTable(myNode, &kademlia.Network)
	// kademlia.Network.TestChannel <- ("My ID : " + myID.String())

	go kademlia.RepublishMyData()
	go kademlia.RepublishData()
	go kademlia.PurgeData()
	conn := kademlia.Network.Listen(myIP)
	kademlia.Network.Conn = conn
	go kademlia.Network.HandleConnection()

	/*	--- Bootstrap Procedure ---	*/
	if bootStrapIP != "" {
		// fmt.Println("GOT bootstrap ID")
		bootStrapNode := NewNode(bootStrapID, bootStrapIP)
		kademlia.Network.SendPingMessage(&kademlia.RoutingTable.me, &bootStrapNode)

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
				go kademlia.LookupNode(&kademlia.RoutingTable.me.ID, queriedNodes, nodeCandidates, 0)

				secondConfirm := <-kademlia.Network.MsgChannel
				if secondConfirm.Command == cmd_find_node_returned {
					var nodeList []Node
					err := json.Unmarshal(secondConfirm.Data, &nodeList)
					checkError(err)
					kademlia.findNodeReturn(&secondConfirm.SenderNode, nodeList)
					kademlia.RoutingTable.RefreshBuckets()
					nodesToPing := kademlia.RoutingTable.FindClosestNodes(&kademlia.RoutingTable.me.ID, kademlia.RoutingTable.GetSize())
					kademlia.RoutingTable.PingNodes(nodesToPing) //pings all known nodes due to above
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

// TODO: FIX THE PATH OS INDEPENDENT
const DOWNLOAD_PATH = "./Downloads"

func (kademlia *Kademlia) channelReader() {
	for {
		//halts here waiting for a command.
		msg := <-kademlia.Network.MsgChannel
		kademlia.RoutingTable.AddNode(msg.SenderNode)
		switch msg.Command {
		case cmd_ping_ack:
			kademlia.RoutingTable.pingedNodes[msg.SenderNode.ID] = false //node has returned ping request.

		case cmd_ping:
			kademlia.Network.SendPingAck(&kademlia.RoutingTable.me, &msg.SenderNode)

		case cmd_store:
			fmt.Println("GOT " + cmd_store)
			var purgeInfo PurgeInformation
			err := json.Unmarshal(msg.Data, &purgeInfo)
			checkError(err)

			// TODO: Add call to own server to establish tcp conn and get the actual file
			select {
			case kademlia.ServerChannel <- NewHandle(CMD_RETRIEVE_FILE, purgeInfo, msg.SenderNode.Address):
				fmt.Println("Sent message to server to get a file")
			case <-time.After(time.Second * 1):
				fmt.Println("Could not deliver retrieve message to server, not listening")
			}
			//QUESTION: Should server handle the below store or kademlia?
			//kademlia.Store(string(msg.Data), DOWNLOAD_PATH, false)

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
			fmt.Println("Found a node that holds the file!")
			//Some node has returned the value that THIS node requested
			select {
			case kademlia.ServerChannel <- NewHandle(CMD_FOUND_FILE, PurgeInformation{Key:string(msg.Data)}, msg.SenderNode.Address):
				fmt.Println("Msg delivered to server")
			case <-time.After(time.Second * 1):
				fmt.Println("Msg could not be delivered to server, server not listening..")
			}

			kademlia.foundHashes[string(msg.Data)] = msg.SenderNode
			// QUESTION: is the below line needed? or handled by server?
			//kademlia.Store(string(msg.Data), DOWNLOAD_PATH, false)

		default:
			fmt.Println("GOT DEFAULT")
		}
		tmpString := ""

		for key, value := range kademlia.files {
			tmpString += "Key: " + key + "Value: " + string(value)
		}

		// kademlia.WriteToFile(kademlia.RoutingTable.me.ID.String()+".txt", []byte(tmpString))
		//kademlia.WriteToFile(kademlia.RoutingTable.me.ID.String()+".txt", []byte(kademlia.RoutingTable.GetRoutingTable()))
	}
}

func (kademlia *Kademlia) PublishData(purgeInfo PurgeInformation, path string, myFile bool) {

	if purgeInfo.Key == "" {
		fmt.Println("In PublishData: Tried to publish empty hash")
		return
	}
	go kademlia.LookupNode(NewKademliaID(purgeInfo.Key), make(map[string]bool), NodeCandidates{}, 0)
	closestNodes := kademlia.RoutingTable.FindClosestNodes(NewKademliaID(purgeInfo.Key), k)
	for i := 0; i < len(closestNodes); i++ {
		marshPurgeInfo, _ := json.Marshal(purgeInfo)
		kademlia.Network.SendStoreMessage(&kademlia.RoutingTable.me, &closestNodes[i], marshPurgeInfo)
	}
	kademlia.Store(purgeInfo, path, myFile)
}

func (kademlia *Kademlia) Get(hash string) string {
	return kademlia.files[hash]
}

func (kademlia *Kademlia) PrintFilesMap() {
	fmt.Println("PRINTING THE KADEMLIA FILES MAP:")
	for key, value := range kademlia.files {
		fmt.Println("Key: ", key, " Value: ", value)
	}
}

func (kademlia *Kademlia) Pin(hash string) {
	tmp := kademlia.Datainfo.PurgeInfos[hash]
	tmp.Pinned = true
	kademlia.Datainfo.PurgeInfos[hash] = tmp
	kademlia.RepublishMyDataOnce()
}

func (kademlia *Kademlia) UnPin(hash string) {
	tmp := kademlia.Datainfo.PurgeInfos[hash]
	tmp.Pinned = false
	kademlia.Datainfo.PurgeInfos[hash] = tmp
	kademlia.RepublishMyDataOnce()
}

func (kademlia *Kademlia) Store(purgeInfo PurgeInformation, path string, me bool) {
	//hash := HashStr(fileName)
	if existingPI, exists := kademlia.Datainfo.PurgeInfos[purgeInfo.Key]; exists {
		fmt.Println("\n YES IT DOES ALREADY EXIST, TIME TO LIVE: ", purgeInfo.TimeToLive, "Pinned: ", purgeInfo.Pinned, "\n")
		existingPI.TimeToLive = purgeInfo.TimeToLive
		existingPI.Pinned = purgeInfo.Pinned
		existingPI.LastPublished = time.Now()
		kademlia.SetPurgeStamp(&existingPI)
		kademlia.Datainfo.PurgeInfos[purgeInfo.Key] = existingPI
	} else {
		fmt.Println("\n NO IT DOES NOT EXIST, TIME TO LIVE: ", purgeInfo.TimeToLive, "\n")
		kademlia.files[purgeInfo.Key] = path
		purgeInfo.LastPublished = time.Now()
		kademlia.SetPurgeStamp(&purgeInfo)
		kademlia.Datainfo.PurgeInfos[purgeInfo.Key] = purgeInfo
	}


	// If the file belongs to this node originally
	if me {
		// Check if the key is already added to myKeys map
		if kademlia.Datainfo.MyKeys[purgeInfo.Key] {
			return
		}
		kademlia.Datainfo.PurgeInfos[purgeInfo.Key] = purgeInfo
		kademlia.Datainfo.MyKeys[purgeInfo.Key] = true
		return
	}

	// If the purgeinformation already exists, update the purgestamp
	// if val, exists := kademlia.Datainfo.PurgeInfos[purgeInfo.Key]; exists {
	// 	kademlia.SetPurgeStamp(&val)
	// 	return
	// }
	//
	// kademlia.SetPurgeStamp(&purgeInfo)
	// kademlia.Datainfo.PurgeInfos[purgeInfo.Key] = purgeInfo

}

func (kademlia *Kademlia) WriteToFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("WRITE ERROR: ", err)
	}
}
