package api

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io"
	"path/filepath"
	"mime"
	"Mobile_D7024E/d7024e"
)

func GetFile(purgeInfo d7024e.PurgeInformation, ip string){
	log.Println("GetFile()")
	log.Println("default_dir = " + default_dir)

	url := "http://"+ip+"/download/"+purgeInfo.Key
	log.Println("     url: " + url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("http.Get error: ", err)
	}

	content := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content)
	if err != nil {
		log.Fatal("mime parse error: ", err)
	}
	filename := params["filename"]

	log.Println("default_dir = " + default_dir)
	path := filepath.Join(default_dir, filename)

	log.Println("path = " + path)

	//If file already exists: return
	if _, err := os.Stat(path); !os.IsNotExist(err) {
 		log.Println("File already exists, will not collect a copy")
		if err != nil {
			log.Println("Error occured in GetFile:", err)
		}
 		return
	}


	out, err := os.Create(path)
	if err != nil  {
		log.Fatal("Create path error: ", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	defer resp.Body.Close()
	if err != nil {
    	log.Fatal("io.Copy error: ", err)
	}
    fmt.Println("Downloaded to "+path)
    kademlia.Store(purgeInfo, path, false)
}

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println("Error in DeleteFile: ", err)
		log.Println("The specified path in DeleteFile is: ", path)
	}
}
