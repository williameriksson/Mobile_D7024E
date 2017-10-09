package d7024e

import (
	"time"
)

//import "fmt"

const bucketSize = 20

type RoutingTable struct {
	me          Node
	buckets     [IDLength * 8]*bucket
	network     *Network
	pingedNodes map[KademliaID]bool
}

func (routingTable *RoutingTable) GetMyAdress() string {
	return routingTable.me.Address
}

func NewRoutingTable(me Node, network *Network) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = NewBucket()
	}
	routingTable.me = me
	routingTable.network = network
	routingTable.pingedNodes = make(map[KademliaID]bool)
	return routingTable
}

func (routingTable *RoutingTable) AddNode(node Node) {
	if !node.ID.Equals(&routingTable.me.ID) {
		bucketIndex := routingTable.GetBucketIndex(&node.ID)
		bucket := routingTable.buckets[bucketIndex]
		if bucket.AddNode(node) {
			//a new node was added or old node updated. all is fine
		} else {
			//a new node was denied due to full bucket, check the bucket!
			go routingTable.addNodeFullBucket(node)
		}
	}
}

//starts asynchronously when a bucket is full to check if node can be added.
func (routingTable *RoutingTable) addNodeFullBucket(node Node) {
	bucketIndex := routingTable.GetBucketIndex(&node.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddToQueue(&node) //adds the new node to the "waitqueue"

	leastNode := []Node{bucket.list.Back().Value.(Node)} //the least recently seen node.
	routingTable.CheckAlive(leastNode)

	bucket.PopQueue()
	newNode := bucket.PopQueue()
	if bucket.Len() < k {
		bucket.AddNode(newNode)
	}
	//at least one node was removed from bucket in question
}

func (routingTable *RoutingTable) RemoveNode(node *Node) {
	if !node.ID.Equals(&routingTable.me.ID) {
		bucketIndex := routingTable.GetBucketIndex(&node.ID)
		bucket := routingTable.buckets[bucketIndex]
		bucket.RemoveNode(node)
	}
}

func (routingTable *RoutingTable) GetBucketSize(bucketIndex int) int {
	return routingTable.buckets[bucketIndex].Len()
}

func (routingTable *RoutingTable) FindClosestNodes(target *KademliaID, count int) []Node {
	var candidates NodeCandidates
	bucketIndex := routingTable.GetBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetNodeAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetNodeAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetNodeAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}
	//fmt.Println(candidates.GetNodes(count))
	// candidates.Print()
	return candidates.GetNodes(count)
}

func (routingTable *RoutingTable) RefreshBuckets() {
	myIndex := routingTable.GetBucketIndex(&routingTable.me.ID)
	//for the buckets less than "me"
	for i := myIndex; i >= 0; i-- {
		if routingTable.buckets[i].Len() < 1 {
			kadID := routingTable.GetRandomIDInBucket(i)
			receiverNode := routingTable.FindClosestNodes(kadID, 1)
			routingTable.network.SendFindNodeMessage(&routingTable.me, &receiverNode[0], kadID)
		}
	}
	//for the buckets more than "me"
	for j := myIndex; j < (IDLength * 8); j++ {
		if routingTable.buckets[j].Len() < 1 {
			kadID := routingTable.GetRandomIDInBucket(j)
			receiverNode := routingTable.FindClosestNodes(kadID, 1)
			routingTable.network.SendFindNodeMessage(&routingTable.me, &receiverNode[0], kadID)
		}
	}
}

func (routingTable *RoutingTable) PingNodes(nodeList []Node) {
	nodes := nodeList
	for i := 0; i < len(nodes); i++ {
		routingTable.network.SendPingMessage(&routingTable.me, &nodes[i])
	}
}

// Checks the provided nodelist if they are still reachable,
// if not (within timeout limit) then they are removed from RoutingTable
func (routingTable *RoutingTable) CheckAlive(nodesToCheck []Node) {
	for i := 0; i < len(nodesToCheck); i++ {
		routingTable.pingedNodes[nodesToCheck[i].ID] = false //set the node as not returned ping yet
	}
	routingTable.PingNodes(nodesToCheck)
	// time.After(timeOutTime)
	time.Sleep(timeOutTime) //wait for returns
	// fmt.Println("nodes to check " + strconv.Itoa(len(nodesToCheck)))
	for i := 0; i < len(nodesToCheck); i++ {
		if routingTable.pingedNodes[nodesToCheck[i].ID] == false {
			routingTable.RemoveNode(&nodesToCheck[i])
		}
	}
}

func (routingTable *RoutingTable) GetBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(&routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}

// func (routingTable *RoutingTable) GetRandomIDInBucket(bucketIndex int) *KademliaID {
// 	kek := routingTable.me.ID.String()
// 	meData, err := hex.DecodeString(kek)
// 	checkError(err)
// 	// fmt.Printf("%X\n", meData)
// 	mask := meData[(bucketIndex/8)] & (1 << (7 - (uint(bucketIndex) % 8)))
// 	// fmt.Printf("%d\n", mask)
// 	if mask == 0 {
// 		meData[(bucketIndex / 8)] |= (1 << (7 - (uint(bucketIndex) % 8)))
// 	} else {
// 		meData[(bucketIndex / 8)] &^= (1 << (7 - (uint(bucketIndex) % 8)))
// 	}
//
// 	str := hex.EncodeToString(meData)
// 	kdID := NewKademliaID(str)
// 	return kdID
// }

func (routingTable *RoutingTable) GetRandomIDInBucket(bucketIndex int) *KademliaID {
	myID := routingTable.me.ID
	randomID := NewRandomKademliaID()
	finalID := NewKademliaID("0000000000000000000000000000000000000000")
	for i := 0; i < IDLength; i++ {
		if i < bucketIndex/8 {
			finalID[i] = myID[i]
		} else if i == bucketIndex/8 {
			finalID[i] = myID[i] ^ (1 << uint(7-(bucketIndex%8)))
			finalID[i] &= (255 << uint(7-(bucketIndex%8)))
			finalID[i] |= randomID[i] >> uint(1+bucketIndex%8)
		} else {
			finalID[i] = randomID[i]
		}
	}
	return finalID
}

func (routingTable *RoutingTable) GetSize() int {
	var size int
	for i := 0; i < (IDLength * 8); i++ {
		size += routingTable.buckets[i].list.Len()
	}
	return size
}

// Returns a string with all nodeID entries of the routingtable,
// for testing purposes
func (routingTable *RoutingTable) GetRoutingTable() string {
	tempString := ""
	for i := 0; i < len(routingTable.buckets); i++ {
		for e := routingTable.buckets[i].list.Front(); e != nil; e = e.Next() {
			nodeID := e.Value.(Node).ID
			tempString += nodeID.String() + "\n"
		}
	}
	return tempString
}

func (routingTable *RoutingTable) GetMyID() string {
	return routingTable.me.ID.String()
}

func (routingTable *RoutingTable) GetMyIP() string {
	return routingTable.me.Address
}
