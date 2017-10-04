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

const default_dir string = "C:/Users/David/go/src/Mobile_D7024E/files/"

func GetFile(hash string, ip string){
	resp, err := http.Get(ip+"/download/"+hash)
	if err != nil {
		log.Fatal(err)
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
    kademlia.Set(hash, path)
}