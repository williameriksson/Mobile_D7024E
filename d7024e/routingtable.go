package d7024e

import (
	"encoding/hex"
	"fmt"
	"strconv"
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
		//check if buckets are full (i.e k amount of nodes)
		if routingTable.GetBucketSize(bucketIndex) >= k {
			go routingTable.addNodeFullBucket(node)
		} else {
			bucket.AddNode(node)
		}
	}
}

//starts asynchronously when a bucket is full to check if node can be added.
func (routingTable *RoutingTable) addNodeFullBucket(node Node) {
	bucketIndex := routingTable.GetBucketIndex(&node.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddToQueue(&node)

	nodes := bucket.GetNodelist()
	routingTable.CheckAlive(nodes)

	for routingTable.GetBucketSize(bucketIndex) < k && bucket.queue.Len() != 0 {
		//at least one node was removed from bucket in question
		bucket.AddNode(bucket.PopQueue())
	}
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
	fmt.Println("nodes to check " + strconv.Itoa(len(nodesToCheck)))
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

func (routingTable *RoutingTable) GetRandomIDInBucket(bucketIndex int) *KademliaID {
	kek := routingTable.me.ID.String()
	meData, err := hex.DecodeString(kek)
	checkError(err)
	// fmt.Printf("%X\n", meData)
	mask := meData[(bucketIndex/8)] & (1 << (7 - (uint(bucketIndex) % 8)))
	// fmt.Printf("%d\n", mask)
	if mask == 0 {
		meData[(bucketIndex / 8)] |= (1 << (7 - (uint(bucketIndex) % 8)))
	} else {
		meData[(bucketIndex / 8)] &^= (1 << (7 - (uint(bucketIndex) % 8)))
	}
	// fmt.Printf("%X\n", meData)
	str := hex.EncodeToString(meData)
	kdID := NewKademliaID(str)
	return kdID
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
