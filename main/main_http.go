package main

import (
	"log"
    "Mobile_D7024E/api"
    "Mobile_D7024E/d7024e"
    "os"
    "fmt"
    "time"
)

// build: go build -o main.exe main.go

var addr string = "127.0.0.1:8100"
var bootstrap string = ""

func main() {
    cmds := os.Args
    //fmt.Println(cmds[1])
    /*cmd := cmds[1]
    switch cmd {
    case "-c":
        kademlia := d7024e.NewKademlia()
        bootstrap := ""
        if len(cmds) == 4 {
            bootstrap = cmds[3]
        }
        kademlia.Run(bootstrap, cmds[2])
    }*/
    switch len(cmds) {
    case 3:
        bootstrap = cmds[2]
        fallthrough
    case 2:
        addr = cmds[1]
        fmt.Println(bootstrap)
        fmt.Println(addr)
    case 1:
    default:
        log.Fatal("Too many arguments.")
    }
    kademlia := d7024e.NewKademlia()
    go kademlia.Run(bootstrap, addr)
    time.Sleep(500*time.Millisecond)
    api.StartServer(kademlia)
}
