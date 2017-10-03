package main

import (
	//"log"
    "Mobile_D7024E/api"
    "Mobile_D7024E/d7024e"
)

// build: go build -o main.exe main.go

func main() {
    kademlia := d7024e.NewKademlia()
    api.StartServer(kademlia)
    kademlia.Run("", "127.0.0.1:8888")
}
