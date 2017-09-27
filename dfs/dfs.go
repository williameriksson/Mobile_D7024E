package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io/ioutil"
)

// Compile command (For Windows): GOOS=windows GOARCH=amd64 go build -o dfs.exe dfs.go

const addr string = "http://127.0.0.1:8080"
const port string = "8080"

func main() {
	cmds := os.Args
	if len(cmds) != 3 {
		log.Fatal("Need exactly 3 arguments.")
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
	default:
		log.Fatal("Unknown argument ", cmds[1])
	}
	//fmt.Println(file)
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