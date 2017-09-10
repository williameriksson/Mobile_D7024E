package main

import (
	"Mobile_D7024E/d7024e"
	"fmt"
)

func main() {
	kademlia := d7024e.Kademlia{}
	go kademlia.Run("", "127.0.0.1:8000")

	kademlia2 := d7024e.Kademlia{}
	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8022")
	for {
		select {
		case msg := <-kademlia2.Network.TestChannel:
			 fmt.Println("received message kad2: ", msg)
	 		default:
	 }
		select {
		case msg := <-kademlia.Network.TestChannel:
			fmt.Println("received message kad: ", msg)
		 default:
	}
	}
}
