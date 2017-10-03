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
}

func (kademlia *Kademlia) PrintHashTable() {
  for key, value := range kademlia.files {
    fmt.Println("Key:", key, "Value:", value)
  }
}

func (kademlia *Kademlia) FindValue(senderNode *Node, hash *KademliaID) {
  // If THIS node has the value, return it
  for key, value := range kademlia.files {
    fmt.Println("Key:", key, "Value:", value)
  }
  fmt.Println("LOOKING FOR THIS HASH: ", strings.TrimSpace(strings.ToLower(hash.String())))
  if val, ok := kademlia.files[strings.TrimSpace(strings.ToLower(hash.String()))]; ok {
    kademlia.Network.SendReturnDataMessage(senderNode, val)
    fmt.Printf("Yes, the value is %x \n", val)
  } else {
    nodeList := kademlia.RoutingTable.FindClosestNodes(hash, k)
    fmt.Println("DID NOT FIND VALUE, THE NODELIST IS: ", nodeList)
    kademlia.Network.SendReturnFindDataMessage(senderNode, nodeList)
  }
}


// If THIS node wants to find a value, it shall call this function. ex:
// kademlia.LookupValue(hash, make(map[string]bool), NodeCandidates{}, 0)
func (kademlia *Kademlia) LookupValue(hash *KademliaID, queriedNodes map[string]bool, prevBestNodes NodeCandidates, recCount int) {
	kademlia.LookupValueCount = 0
	kademlia.returnedValueNodes = prevBestNodes
  fmt.Println("THE VALUE IS: ", kademlia.returnedValue)

	closestNodes := kademlia.RoutingTable.FindClosestNodes(hash, k)
	//fmt.Println("closest len " + strconv.Itoa(len(closestNodes)))
	for i := 0; i < alpha && i < len(closestNodes); i++ {
		if queriedNodes[hash.String()] == false {
			queriedNodes[hash.String()] = true
      fmt.Println("!!!!!SENDING THE FIND DATA MESSAGE!!!!!")
      kademlia.Network.SendFindDataMessage(&closestNodes[i], hash)
		}
	}

	timeout := false
	timeStamp := time.Now()

	for (kademlia.LookupValueCount < 1) && !timeout {
		//busy waiting for at least one RETURN_FIND_NODE
		if time.Now().Sub(timeStamp) > timeOutTime {
			timeout = true
		}
	}

  for i := 0; i < len(kademlia.returnedValueNodes.nodes); i++ {
    kademlia.returnedValueNodes.nodes[i].CalcDistance(hash)
  }
  kademlia.returnedValueNodes.Sort()

	if !timeout || timeout {
    if len(kademlia.returnedValue) != 0 {
      // We got the value, now we cache it in the closest node to the target that didn't have it
      // How do we make sure this isn't the node that we got it from? Maybe we have to save
      // the id of the node that sent the data and remove it from consideration. For now we just take
      // the node with the closest ID to the key
      if (len(kademlia.returnedValueNodes.nodes) > 0) {
        kademlia.Network.SendStoreMessage(&kademlia.returnedValueNodes.nodes[0], kademlia.returnedValue)
      }

      fmt.Println("Value found, let's do something with it.", kademlia.returnedValue)

      kademlia.returnedValue = nil
      return
    }
    //fmt.Println("RETURNEDVALUENODES: ", kademlia.returnedValueNodes)


    var length int
    if length = k; kademlia.returnedValueNodes.Len() < k {
      length = kademlia.returnedValueNodes.Len()
    }


		bestNodes := NodeCandidates{nodes: kademlia.returnedValueNodes.GetNodes(length)}
		if recCount == k {
			//did NOT find the data after k attempts
		} else {
			//did NOT find data, continue search
			kademlia.LookupValue(hash, queriedNodes, bestNodes, (recCount + 1))
		}
	}




}
