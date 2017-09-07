package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	me          *Node
	msgChannel  chan Message
	testChannel chan string
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

func (network *Network) Listen(myIP string) {
	// fmt.Println(myIP)
	serverAddr, err := net.ResolveUDPAddr("udp", myIP)
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	// fmt.Println(conn)
	go network.handleConnection(conn)
	// go network.test()
	for {

	}
}

func (network *Network) handleConnection(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	var msg Message

	defer conn.Close()
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		checkError(err)
		err = json.Unmarshal(buf[0:n], &msg)
		checkError(err)

		switch msg.Command {
		case "PING_ACK":
			fmt.Println("GOT PING_ACK")
			network.PingAck(&msg)
		case "PING":
			network.testChannel <- network.me.Address
			fmt.Println("GOT PING")
		case "STORE":
			fmt.Println("GOT STORE")
		case "FIND_NODE":
			fmt.Println("GOT FIND_NODE")
		case "FIND_VALUE":
			fmt.Println("FIND_VALUE")
		default:
			fmt.Println("GOT DEFAULT", n, addr, msg)
		}
	}
}

func (network *Network) sendMessage(receiverNode *Node, msg *Message) {
	ServerAddr, err := net.ResolveUDPAddr("udp", receiverNode.Address) // take from routingtable.me.address
	checkError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", network.me.Address)
	checkError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	checkError(err)

	defer Conn.Close()

	marshMsg, err := json.Marshal(msg)
	checkError(err)
	Conn.Write(marshMsg)
}

func (network *Network) SendPingMessage(receiverNode *Node) {
	go network.sendMessage(receiverNode, &Message{Command: "PING", SenderNode: network.me})
}

func (network *Network) PingAck(msg *Message) {
	network.msgChannel <- *msg
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
