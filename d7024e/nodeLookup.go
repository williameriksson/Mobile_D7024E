package d7024e

import (
	"time"
)

func (kademlia *Kademlia) LookupNode(target *KademliaID, queriedNodes map[string]bool, prevBestNodes NodeCandidates, recCount int) {
	// Should be run when looking for a node (not during bootstrap though)
	lookUpCountMutex.Lock()
	kademlia.LookupCount = 0
	lookUpCountMutex.Unlock()

	returnedNodesMutex.Lock()
	kademlia.returnedNodes = prevBestNodes
	returnedNodesMutex.Unlock()

	// fmt.Println("LookupNode running")
	closestNodes := kademlia.RoutingTable.FindClosestNodes(target, k)
	for i := 0; i < alpha && i < len(closestNodes); i++ {
		if queriedNodes[target.String()] == false {
			queriedNodes[target.String()] = true
			kademlia.Network.SendFindNodeMessage(&kademlia.RoutingTable.me, &closestNodes[i], target)
		}
	}

	timeout := false
	select {
	case <-kademlia.timeoutChannel:
		timeout = false
	case <-time.After(time.Second * 1):
		timeout = true
	}

	if !timeout {
		returnedNodesMutex.Lock()
		kademlia.returnedNodes.Sort() //what does this do
		count := kademlia.returnedNodes.Len()
		if count > k {
			count = k
		}
		bestNodes := NodeCandidates{nodes: kademlia.returnedNodes.GetNodes(count)}
		returnedNodesMutex.Unlock()

		if count == 0 {
			kademlia.LookupNode(target, queriedNodes, prevBestNodes, (recCount + 1))
		} else if bestNodes.nodes[0].ID.String() == target.String() {
			//first node IS target means we DID find it this run.
		} else if recCount == k {
			//did NOT find node after k attempts
		} else {
			//did NOT find node, continue search
			kademlia.LookupNode(target, queriedNodes, bestNodes, (recCount + 1))
		}
	} else {
		kademlia.LookupNode(target, queriedNodes, prevBestNodes, (recCount + 1))
	}
}

func (kademlia *Kademlia) findNode(senderNode *Node, kID *KademliaID) {
	nodeList := kademlia.RoutingTable.FindClosestNodes(kID, k)
	kademlia.Network.SendReturnFindNodeMessage(&kademlia.RoutingTable.me, senderNode, nodeList)
}

func (kademlia *Kademlia) findNodeReturn(senderNode *Node, nodeList []Node) {
	returnedNodesMutex.Lock()
	kademlia.returnedNodes.Append(nodeList)
	returnedNodesMutex.Unlock()

	lookUpCountMutex.Lock()
	kademlia.LookupCount++
	lookUpCountMutex.Unlock()

	//adds all the returned nodes to the RoutingTable
	for i := 0; i < len(nodeList); i++ {
		kademlia.RoutingTable.AddNode(nodeList[i])
	}
	select {
	case kademlia.timeoutChannel <- true:
	default:
	}
}
