package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"io"
	"strconv"
	"time"
	"path/filepath"
	"mime"
)

// Compile command: go build -o dfs.exe dfs.go

const addr string = "http://127.0.0.1:8080"
const dir string = "C:/Users/David/go/src/Mobile_D7024E/files/"

func main() {
	cmds := os.Args
	if len(cmds) < 3 {
		log.Fatal("Too few arguments.")
	}

	switch cmds[1] {
	case "store":
		path, err := filepath.Abs(cmds[2])
		if err != nil {
			fmt.Println(err)
		}

		// Check if file exists
		if file, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatal("File doesn't exist.")
		} else if file.Mode().IsDir(){
			// If "file" is a directory it throws fatal error
			log.Fatal("Not a file.")
		}
		
		val, err := url.ParseQuery("path="+path)
		if err != nil {
			fmt.Println(err)
		}
		post(addr+"/store", val)
		
		//fmt.Println(path)
	case "cat":
		resp, err := http.Get(addr+"/cat/"+cmds[2])
		if err != nil {
			log.Fatal(err)
		}

		text, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
	    	log.Fatal(err)
		}
	    fmt.Print(string(text))
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
		//get(addr+"/addnode/127.0.0.1:"+strconv.Itoa(port))

		nr, _ := strconv.Atoi(cmds[2])
		for i := 1; i < nr; i++ {
			get(addr+"/addnode/127.0.0.1:"+strconv.Itoa(port+i)+"?bootstrap=127.0.0.1:"+strconv.Itoa(port))
			time.Sleep(time.Millisecond * 1000)
		}
	case "download":
		resp, err := http.Get(addr+"/download/"+cmds[2])
		if err != nil {
			log.Fatal(err)
		}

		content := resp.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(content)
		if err != nil {
			log.Fatal(err)
		}
		filename := params["filename"]
		
		out, err := os.Create(dir + filename)
		if err != nil  {
			log.Fatal(err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		defer resp.Body.Close()
		if err != nil {
	    	log.Fatal(err)
		}
	    fmt.Println("Downloaded to "+dir+filename)
	default:
		log.Fatal("Unknown argument ", cmds[1])
	}
}

func post(url string, data url.Values) {
	resp, err := http.PostForm(url, data)

	text, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
    	log.Fatal(err)
	}
    fmt.Println(string(text))
}

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	text, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
    	log.Fatal(err)
	}
    fmt.Println(string(text))
}