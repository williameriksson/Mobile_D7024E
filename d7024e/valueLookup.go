package d7024e

import (
  "fmt"
  "time"
  "strings"
)

func (kademlia *Kademlia) FindValueReturn(senderNode *Node, nodeList []Node) {
	kademlia.returnedValueNodes.Append(nodeList)
	kademlia.LookupValueCount++

	//adds all the returned nodes to the RoutingTable
	for i := 0; i < len(nodeList); i++ {
		kademlia.RoutingTable.AddNode(nodeList[i])
	}
	select {
	case kademlia.valueTimeoutChan <- true:
	default:
	}
}

func (kademlia *Kademlia) PrintHashTable() {
	for key, value := range kademlia.files {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func (kademlia *Kademlia) FindValue(senderNode *Node, hash *KademliaID) {
  // If THIS node has the value, echo the hash back to the requester
  for key, value := range kademlia.files {
    fmt.Println("Key:", key, "Value:", value)
  }
  fmt.Println("LOOKING FOR THIS HASH: ", strings.TrimSpace(strings.ToLower(hash.String())))
  if val, ok := kademlia.files[strings.TrimSpace(strings.ToLower(hash.String()))]; ok {
    kademlia.Network.SendReturnDataMessage(&kademlia.RoutingTable.me, senderNode, []byte(hash.String()))
    fmt.Printf("Yes, the value is %v \n", val)
  } else {
    nodeList := kademlia.RoutingTable.FindClosestNodes(hash, k)
    // fmt.Println("DID NOT FIND VALUE, THE NODELIST IS: ", nodeList)
    kademlia.Network.SendReturnFindDataMessage(&kademlia.RoutingTable.me, senderNode, nodeList)
  }
}

func (kademlia *Kademlia) LookupValue(hash string) {
  hashId := NewKademliaID(hash)
  kademlia.lookupValue(hashId, make(map[string]bool), NodeCandidates{}, 0)
}

// If THIS node wants to find a value, it shall call this function. ex:
// kademlia.LookupValue(hash, make(map[string]bool), NodeCandidates{}, 0)
func (kademlia *Kademlia) lookupValue(hash *KademliaID, queriedNodes map[string]bool, prevBestNodes NodeCandidates, recCount int) {
	kademlia.LookupValueCount = 0
	kademlia.returnedValueNodes = prevBestNodes

	closestNodes := kademlia.RoutingTable.FindClosestNodes(hash, k)
	//fmt.Println("closest len " + strconv.Itoa(len(closestNodes)))
	for i := 0; i < alpha && i < len(closestNodes); i++ {
		if queriedNodes[hash.String()] == false {
			queriedNodes[hash.String()] = true
      if closestNodes[i].ID != kademlia.RoutingTable.me.ID {
        fmt.Println("!!!!!SENDING THE FIND DATA MESSAGE!!!!!")
        kademlia.Network.SendFindDataMessage(&kademlia.RoutingTable.me, &closestNodes[i], hash)
      }
		}
	}

	timeout := false
	select {
	case <-kademlia.valueTimeoutChan:
		timeout = false
	case <-time.After(time.Second * 1):
		timeout = true
	}

	for i := 0; i < len(kademlia.returnedValueNodes.nodes); i++ {
		kademlia.returnedValueNodes.nodes[i].CalcDistance(hash)
	}
	kademlia.returnedValueNodes.Sort()

	if !timeout {
    if _, exists := kademlia.foundHashes[hash.String()]; exists {
      // We got the value, now we cache it in the closest node to the target that didn't have it
      for _, node := range kademlia.returnedValueNodes.nodes {
          if node.ID != kademlia.foundHashes[hash.String()].ID {
            kademlia.Network.SendStoreMessage(&kademlia.RoutingTable.me, &kademlia.returnedValueNodes.nodes[0], []byte(hash.String()))
            delete(kademlia.foundHashes, hash.String())
            return
          }
      }
      return
    }
    var length int
    if length = k; kademlia.returnedValueNodes.Len() < k {
      length = kademlia.returnedValueNodes.Len()

	    bestNodes := NodeCandidates{nodes: kademlia.returnedValueNodes.GetNodes(length)}
    	if recCount == k {
    		//did NOT find the data after k attempts
    	} else {
    		//did NOT find data, continue search
    		kademlia.lookupValue(hash, queriedNodes, bestNodes, (recCount + 1))
    	}
    }

  } else {
    kademlia.lookupValue(hash, queriedNodes, prevBestNodes, (recCount + 1))
  }

}
