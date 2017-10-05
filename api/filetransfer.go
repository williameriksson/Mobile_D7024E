package api

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io"
	//"path/filepath"
	"mime"
)

func GetFile(hash string, ip string){
	log.Println("GetFile()")
	resp, err := http.Get("http://"+ip+"/download/"+hash)
	if err != nil {
		log.Fatal("http.Get error: "+ err)
	}

	content := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content)
	if err != nil {
		log.Fatal(err)
	}
	filename := params["filename"]
	path := default_dir + filename
	
	out, err := os.Create(path)
	if err != nil  {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	defer resp.Body.Close()
	if err != nil {
    	log.Fatal(err)
	}
    fmt.Println("Downloaded to "+path)
    kademlia.Store(hash, path, false)
}