package main

import (
	"Mobile_D7024E/d7024e"
	// "fmt"
	"time"
	"strconv"
)

func main() {

	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 1000)

	var port int = 8001
	for i := 0; i < 10; i++ {
		go d7024e.NewKademlia().Run("127.0.0.1:8000", "127.0.0.1:" + strconv.Itoa(port + i))
		time.Sleep(time.Millisecond * 1000)
	}




	// kademlia2 := d7024e.NewKademlia()
	// go kademlia2.Run("127.0.0.1:8000", "127.0.0.1:8002")
	//
	// time.Sleep(time.Millisecond * 1000)
	//
	// kademlia3 := d7024e.NewKademlia()
	// go kademlia3.Run("127.0.0.1:8002", "127.0.0.1:8003")
	//
	// time.Sleep(time.Millisecond * 1000)
	//
	// kademlia4 := d7024e.NewKademlia()
	// go kademlia4.Run("127.0.0.1:8003", "127.0.0.1:8004")
	//
	// time.Sleep(time.Millisecond * 1000)
	//
	// kademlia5 := d7024e.NewKademlia()
	// go kademlia5.Run("127.0.0.1:8004", "127.0.0.1:8005")
	//
	// time.Sleep(time.Millisecond * 1000)
	//
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
