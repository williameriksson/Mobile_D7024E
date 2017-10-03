package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	me          *Node
	MsgChannel  chan Message
	TestChannel chan string
	Conn        *net.UDPConn
}

type Message struct {
	Command    string
	SenderNode Node
	Hash       string
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

func (network *Network) SendPingMessage(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_ping, SenderNode: *network.me})
}

func (network *Network) SendPingAck(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_ping_ack, SenderNode: *network.me})
}

func (network *Network) SendFindNodeMessage(receiverNode *Node, kiD *KademliaID) {
	data := []byte(kiD.String())
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_node, SenderNode: *network.me, Data: data})
}

func (network *Network) SendReturnFindNodeMessage(receiverNode *Node, nodeList []Node) {
	data, err := json.Marshal(nodeList)
	checkError(err)
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_node_returned, SenderNode: *network.me, Data: data})
}

func (network *Network) SendFindDataMessage(receiverNode *Node, hash string) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_find_value, SenderNode: *network.me, Hash: hash})
}

func (network *Network) SendStoreMessage(receiverNode *Node, data []byte) {
	go network.sendMessage(receiverNode, &Message{Command: cmd_store, SenderNode: *network.me, Data: data})
}
