package d7024e

import (
	"fmt"
	"sort"
)

type Node struct {
	ID       KademliaID
	Address  string
	distance *KademliaID
}

func NewNode(id *KademliaID, address string) Node {
	return Node{*id, address, nil}
}

func (node *Node) CalcDistance(target *KademliaID) {
	node.distance = node.ID.CalcDistance(target)
}

func (node *Node) Less(otherNode *Node) bool {
	return node.distance.Less(otherNode.distance)
}

func (node *Node) String() string {
	return fmt.Sprintf(`node("%s", "%s")`, node.ID, node.Address)
}

type NodeCandidates struct {
	nodes []Node
}

func (candidates *NodeCandidates) Append(nodes []Node) {
	candidates.nodes = append(candidates.nodes, nodes...)
}

func (candidates *NodeCandidates) GetNodes(count int) []Node {
	return candidates.nodes[:count]
}

func (candidates *NodeCandidates) Sort() {
	sort.Sort(candidates)
}

func (candidates *NodeCandidates) Len() int {
	return len(candidates.nodes)
}

func (candidates *NodeCandidates) Swap(i, j int) {
	candidates.nodes[i], candidates.nodes[j] = candidates.nodes[j], candidates.nodes[i]
}

func (candidates *NodeCandidates) Less(i, j int) bool {
	return candidates.nodes[i].Less(&candidates.nodes[j])
}

func (candidates *NodeCandidates) Print() {
	for _, node := range candidates.nodes {
		fmt.Printf("ID: 0x%X, IP: %v, Distance: 0x%X \n", node.ID, node.Address, *node.distance)
	}
	
}
