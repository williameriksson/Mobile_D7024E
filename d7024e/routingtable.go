package d7024e

//import "fmt"

const bucketSize = 20

type RoutingTable struct {
	me      Node
	buckets [IDLength * 8]*bucket
}

func NewRoutingTable(me Node) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = NewBucket()
	}
	routingTable.me = me
	return routingTable
}

func (routingTable *RoutingTable) AddNode(node Node) {
	if !node.ID.Equals(&routingTable.me.ID) {
		bucketIndex := routingTable.getBucketIndex(&node.ID)
		bucket := routingTable.buckets[bucketIndex]
		bucket.AddNode(node)
	}
}

func (routingTable *RoutingTable) FindClosestNodes(target *KademliaID, count int) []Node {
	var candidates NodeCandidates
	bucketIndex := routingTable.getBucketIndex(target)
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

func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
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
