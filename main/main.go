package main

import (
	"Mobile_D7024E/d7024e"
	//"fmt"
	"strconv"
	"time"
	//"bytes"
)

// build: go build -o main.exe main.go
func main() {


	kademlia1 := d7024e.NewKademlia()
	go kademlia1.Run("", "127.0.0.1:8000")
	time.Sleep(time.Millisecond * 1000)

	// var lastKademlia d7024e.Kademlia
	var port int = 8001
	for i := 0; i < 4; i++ {
		go d7024e.NewKademlia().Run("127.0.0.1:8000", "127.0.0.1:"+strconv.Itoa(port+i))
		time.Sleep(time.Millisecond * 500)
	}

	// data := []byte("Test")
	// kademlia1.PublishData(data)
	//
	// time.Sleep(time.Millisecond * 1500)
	// hash := d7024e.HashData(data)
	// lastKademlia.LookupValue(d7024e.NewKademliaID(hash), make(map[string]bool), d7024e.NodeCandidates{}, 0)
	// kademlia1.LookupValue(d7024e.HashData(data))
	/* End of Davids Test Code */
	for {

	}
}
