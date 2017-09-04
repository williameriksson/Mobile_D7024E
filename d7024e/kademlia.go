package d7024e

import "fmt"

const alpha int = 3

type Kademlia struct {
  test int
  routingTable RoutingTable
}

type Export struct {

}

func LeTest(str string) {
  fmt.Println(str)
}

func (kademlia *Kademlia) LookupNode(target *Node) {
	// TODO

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
