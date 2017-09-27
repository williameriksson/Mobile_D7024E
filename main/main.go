package main

import (
	"Mobile_D7024E/d7024e"
	//"fmt"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"time"
	//"bytes"
)

func main() {

	c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
      <-c
      os.Exit(1)
  }()

	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 1000)

	var lastKademlia d7024e.Kademlia
	var port int = 8001
	for i := 0; i < 2; i++ {
		lastKademlia = *d7024e.NewKademlia()
		go lastKademlia.Run("127.0.0.1:8000", "127.0.0.1:"+strconv.Itoa(port+i))
		time.Sleep(time.Millisecond * 1000)
	}



	data := []byte("Test")
	kademlia1.PublishData(data)
	/* Davids Test Code */

	// n := bytes.IndexByte(data, 0)
	// dataHash := d7024e.NewKademliaID(string(data[:n]))
	// kademlia1.Store(data)
	// kademlia1.PrintHashTable()
	// fmt.Println("THE DATA:", data)
	// fmt.Println("THE HASH:", hash)
	time.Sleep(time.Millisecond * 1500)
	hash := d7024e.HashData(data)
	lastKademlia.LookupValue(d7024e.NewKademliaID(hash), make(map[string]bool), d7024e.NodeCandidates{}, 0)
	// kademlia1.LookupValue(d7024e.HashData(data))
	/* End of Davids Test Code */
	for {
		// select {
		// case msg := <-kademlia1.Network.TestChannel:
		// 	fmt.Println(msg)
		// case msg := <-kademlia2.Network.TestChannel:
		// 	fmt.Println(msg)
		// case msg := <-kademlia3.Network.TestChannel:
		// 	fmt.Println(msg)
		// case msg := <-kademlia4.Network.TestChannel:
		// 	fmt.Println(msg)
		// case msg := <-kademlia5.Network.TestChannel:
		// 	fmt.Println(msg)
		// }
	}
}
