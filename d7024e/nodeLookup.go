package d7024e

import (
	"fmt"
	"time"
)

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
	select {
	case <-kademlia.timeoutChannel:
		timeout = false
	case <-time.After(time.Second * 1):
		timeout = true
	}

	if !timeout {
		// fmt.Println("RETURNEDNODES: ", kademlia.returnedValueNodes)
		kademlia.returnedNodes.Sort() //what does this do
		count := kademlia.returnedNodes.Len()
		if count > k {
			count = k
		}
		bestNodes := NodeCandidates{nodes: kademlia.returnedNodes.GetNodes(count)}
		if bestNodes.nodes[0].ID.String() == target.String() {
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
	select {
	case kademlia.timeoutChannel <- true:
	default:
	}
}
