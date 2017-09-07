package main

import "Mobile_D7024E/d7024e"

func main() {
	kademlia := d7024e.Kademlia{}
	go kademlia.Run("", "127.0.0.1:8023")

	// kademlia2 := d7024e.Kademlia{}
	// go kademlia2.Run("127.0.0.1:8023", "127.0.0.1:8015")
	for {

	}
}
