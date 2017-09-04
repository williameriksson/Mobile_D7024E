package d7024e

import (
	"fmt"
	"net"
	"time"
	"strconv"
)

type Network struct {
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (network *Network) Listen() {
	serverAddr, err := net.ResolveUDPAddr("udp",":8000")
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	fmt.Println("1")
	channel1 := make(chan string)
	go network.handleConnection(conn, channel1)
	go network.test()
	for {
		msg := <-channel1
		fmt.Println(msg)
	}
	fmt.Println("2")
}

func (network *Network) handleConnection(conn *net.UDPConn, channel1 chan string) {
	buf := make([]byte, 1024)
	fmt.Println("3")
	defer conn.Close()
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		checkError(err)
		receivedString := string(buf[0:n])
		fmt.Println(n, addr, receivedString)
		channel1 <- receivedString
	}
	fmt.Println("4")
}

func (network *Network) test() {
		fmt.Println("5")
    ServerAddr,err := net.ResolveUDPAddr("udp",":8000")
    checkError(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", ":8001")
    checkError(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    checkError(err)

    defer Conn.Close()
    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
    }
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
