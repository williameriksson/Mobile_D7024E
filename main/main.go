package main

import (
	"Mobile_D7024E/d7024e"
	// "fmt"
	"strconv"
	"time"
)

func main() {

	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 1000)

	var port int = 8001
	for i := 0; i < 25; i++ {
		go d7024e.NewKademlia().Run("127.0.0.1:8000", "127.0.0.1:"+strconv.Itoa(port+i))
		time.Sleep(time.Millisecond * 500)
	}

	/* Davids Test Code */
	// data := d7024e.HashStr("Teststring")
	// kademlia.Store(data)
	// kademlia.LookupData(d7024e.HashData(data))
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
