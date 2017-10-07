package api

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io"
	//"path/filepath"
	"mime"
	"Mobile_D7024E/d7024e"
)

func GetFile(purgeInfo d7024e.PurgeInformation, ip string){
	log.Println("GetFile()")

	url := "http://"+ip+"/download/"+purgeInfo.Key
	log.Println("     url: " + url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("http.Get error: ", err)
	}

	content := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content)
	if err != nil {
		log.Fatal(err)
	}
	filename := params["filename"]
	path := default_dir + filename

	// If file already exists: return
	if _, err := os.Stat(path); !os.IsNotExist(err) {
 		log.Println(err)
 		return
	}


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
    kademlia.Store(purgeInfo, path, false)
}

func DeleteFile(path string){
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
	}
}
