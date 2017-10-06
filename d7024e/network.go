package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	MsgChannel  chan Message
	TestChannel chan string
	Conn        *net.UDPConn
}

type Message struct {
	Command    string
	SenderNode Node
	Data       []byte
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (network *Network) Listen(myIP string) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", myIP)
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	return conn
}

func (network *Network) HandleConnection() {
	buf := make([]byte, 65507)
	var msg Message

	defer network.Conn.Close()
	for {
		n, _, err := network.Conn.ReadFromUDP(buf)
		checkError(err)
		err = json.Unmarshal(buf[0:n], &msg)
		checkError(err)
		network.MsgChannel <- msg
	}
}

func (network *Network) sendMessage(receiverNode *Node, msg *Message) {
	receiverAddr, err := net.ResolveUDPAddr("udp", receiverNode.Address)
	checkError(err)

	marshMsg, err := json.Marshal(msg)
	checkError(err)
	network.Conn.WriteToUDP(marshMsg, receiverAddr)
}

func (network *Network) SendPingMessage(senderNode *Node, receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_ping, SenderNode: *senderNode})
}

func (network *Network) SendPingAck(senderNode *Node, receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_ping_ack, SenderNode: *senderNode})
}

func (network *Network) SendFindNodeMessage(senderNode *Node, receiverNode *Node, kiD *KademliaID) {
	//data := []byte(kiD.String())
	data, err := json.Marshal(kiD)
	checkError(err)
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_node, SenderNode: *senderNode, Data: data})
}

func (network *Network) SendReturnFindNodeMessage(senderNode *Node, receiverNode *Node, nodeList []Node) {
	data, err := json.Marshal(nodeList)
	checkError(err)
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_node_returned, SenderNode: *senderNode, Data: data})
}

func (network *Network) SendFindDataMessage(senderNode *Node, receiverNode *Node, hash *KademliaID) {
	data, err := json.Marshal(hash)
	checkError(err)
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_value, SenderNode: *senderNode, Data: data})
}

// If you don't hold the requested data, return the closest nodes list.
func (network *Network) SendReturnFindDataMessage(senderNode *Node, receiverNode *Node, nodeList []Node) {
	data, err := json.Marshal(nodeList)
	checkError(err)
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_value_returned, SenderNode: *senderNode, Data: data})
}

// If you hold the requested data, return it.
func (network *Network) SendReturnDataMessage(senderNode *Node, receiverNode *Node, data []byte) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_value_returned, SenderNode: *senderNode, Data: data})
}

func (network *Network) SendStoreMessage(senderNode *Node, receiverNode *Node, data []byte) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_store, SenderNode: *senderNode, Data: data})
}
