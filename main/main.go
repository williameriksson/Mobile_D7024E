package main

import (
	"Mobile_D7024E/d7024e"
	"fmt"
)

func main() {

	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")

	// time.Sleep(time.Millisecond * 1000)

	kademlia2 := d7024e.NewKademlia()
	go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8055")

	// time.Sleep(time.Millisecond * 1000)
	//
	// kademlia3 := d7024e.NewKademlia()
	// go kademlia3.Run("127.0.0.1:8000", "127.0.0.1:8033")

	//
	/* Davids Test Code */
	// data := d7024e.HashStr("Teststring")
	// kademlia.Store(data)
	// kademlia.LookupData(d7024e.HashData(data))
	/* End of Davids Test Code */

	for {
		select {
		case msg := <-kademlia1.Network.TestChannel:
			fmt.Println(msg)
		default:
		}
		select {
		case msg := <-kademlia2.Network.TestChannel:
			fmt.Println(msg)
		default:
		}
		// select {
		// case msg := <-kademlia3.Network.TestChannel:
		// 	fmt.Println(msg)
		// default:
		// }
	}
}
