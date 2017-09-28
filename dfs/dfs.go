package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
)

// Compile command (For Windows): GOOS=windows GOARCH=amd64 go build -o dfs.exe dfs.go

const addr string = "http://127.0.0.1:8080"

func main() {
	cmds := os.Args
	if len(cmds) < 3 {
		log.Fatal("Too few arguments.")
	}

	switch cmds[1] {
	case "store":
		fmt.Println("Store")
	case "cat":
		get(addr+"/cat/"+cmds[2])
	case "pin":
		get(addr+"/pin/"+cmds[2])
	case "unpin":
		get(addr+"/unpin/"+cmds[2])

	/*
	 * dfs addnode nodeIP bootstrapIP
	 * creates new node that connects through the bootstrapIP node (bootstrapIP is optional).
	 */
	case "addnode":
		bootstrap := ""
		if len(cmds) == 4 {
			bootstrap = cmds[3]
		}
		get(addr+"/addnode/"+cmds[2]+"?bootstrap="+bootstrap)

	/* 
	 * dfs populate X
	 * makes a new Kademlia system with X nodes
	 */ 
	case "populate":
		var port int = 8100
		get(addr+"/addnode/127.0.0.1:"+strconv.Itoa(port))

		nr, _ := strconv.Atoi(cmds[2])
		for i := 1; i < nr; i++ {
			get(addr+"/addnode/127.0.0.1:"+strconv.Itoa(port+i)+"?bootstrap=127.0.0.1:"+strconv.Itoa(port))
			time.Sleep(time.Millisecond * 1000)
		}
	default:
		log.Fatal("Unknown argument ", cmds[1])
	}
}

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    	log.Fatal(err)
	}
    fmt.Println(string(text))
}