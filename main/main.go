package main

import (
	"Mobile_D7024E/d7024e"
)

func main() {
	/*kademlia := d7024e.Kademlia{}
	go kademlia.Run("", "127.0.0.1:16667")

	kademlia2 := d7024e.Kademlia{}
	go kademlia2.Run("127.0.0.1:16667", "127.0.0.1:10111")
	for {
	
	}*/
	kademlia := d7024e.NewKademlia()
	data := d7024e.HashStr("Teststring")
	kademlia.Store(data)
	kademlia.LookupData(d7024e.HashData(data))

}
