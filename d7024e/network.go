package d7024e

import (
	"fmt"
	"net"
	"encoding/json"
)

type Network struct {
}

type Message struct {
	Command string
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (network *Network) Listen() {
	serverAddr, err := net.ResolveUDPAddr("udp",":8002")
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	go network.handleConnection(conn)
	go network.test()
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
		case "PING":
			fmt.Println("GOT PING")
		default:
			fmt.Println("GOT DEFAULT", n, addr, msg)
		}
	}

}

func (network *Network) test() {
    ServerAddr,err := net.ResolveUDPAddr("udp",":8002")
    checkError(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", ":8003")
    checkError(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    checkError(err)

    defer Conn.Close()

		msg := Message{"PING"}
		marshMsg, err := json.Marshal(msg)
		checkError(err)
		Conn.Write(marshMsg)

}

/*
func (network *Network) SendPingMessage(node *Node) {
	// TODO
}

func (network *Network) SendFindNodeMessage(node *Node) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
*/
