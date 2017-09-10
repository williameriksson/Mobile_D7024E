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
	Conn				*net.UDPConn
}

type Message struct {
	Command    string
	SenderNode *Node
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (network *Network) Listen(myIP string) *net.UDPConn{
	// fmt.Println(myIP)
	serverAddr, err := net.ResolveUDPAddr("udp", myIP)
	conn, err := net.ListenUDP("udp", serverAddr)
	//network.Conn = conn
	checkError(err)
	return conn
	//go network.handleConnection(conn)
	// go network.test()
}

func (network *Network) HandleConnection() {
	buf := make([]byte, 1024)
	var msg Message

	defer network.Conn.Close()
	for {
		n, addr, err := network.Conn.ReadFromUDP(buf)
		checkError(err)
		err = json.Unmarshal(buf[0:n], &msg)
		checkError(err)

		switch msg.Command {
		case "PING_ACK":
			fmt.Println("GOT PING_ACK")
			go network.PingAck(&msg)
		case "PING":
			network.TestChannel <- network.me.Address
			network.SendPingAck(msg.SenderNode)
			fmt.Println("GOT PING")
		case "STORE":
			fmt.Println("GOT STORE")
		case "FIND_NODE":
			fmt.Println("GOT FIND_NODE")
		case "FIND_VALUE":
			fmt.Println("FIND_VALUE")
		default:
			fmt.Println("GOT DEFAULT", n, addr, &msg)
		}
	}
}

func (network *Network) sendMessage(receiverNode *Node, msg *Message) {
	ServerAddr, err := net.ResolveUDPAddr("udp", receiverNode.Address)
	checkError(err)

	marshMsg, err := json.Marshal(msg)
	checkError(err)
	network.Conn.WriteToUDP(marshMsg, ServerAddr)
}

func (network *Network) SendPingMessage(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: "PING", SenderNode: network.me})
}

func (network *Network) SendPingAck(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: "PING_ACK", SenderNode: network.me})
}

func (network *Network) PingAck(msg *Message) {
	network.TestChannel <- "Got PING_ACK"
	network.MsgChannel <- *msg
}

func (network *Network) SendFindNodeMessage(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: "FIND_NODE", SenderNode: network.me})
}

/*
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
*/
