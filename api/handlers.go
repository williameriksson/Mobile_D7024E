package api

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	"net/url"
	//"strconv"
	"path/filepath"
	"log"
	"os"
	"io"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"Mobile_D7024E/d7024e"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kademlia API\n\nLocalhost:\n     GET:  /pin/{hash}\n     GET:  /unpin/{hash}\n     GET:  /addnode/{addr}?boostrap={bootstrap_addr}\n     POST: /store/\nPublic:\n     GET:  /cat/{hash}\n     GET:  /download/{hash}")
}

func Cat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	log.Println("hash:"+hash)
	fp := kademlia.Get(hash)
	log.Println("fp:"+fp)
	/*w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))*/
	// Create the file
	out, err := os.Open(fp)
	defer out.Close()
	if err != nil  {
		log.Println("I don't have file: ", err)

		// Search Kademlia Network for file
		go kademlia.LookupValue(hash)

		var ip string

		select {

		// Wait on response from kademlia
		case ip = <- res:
		// Timeout
		case <- time.After(time.Second * 2):
			log.Println("Timeout in Cat() function.")
			fmt.Fprint(w, "Error, Timeout in Cat() function.")
			return
		}

		// Change from kademlia port to API port
		ip = convertIP(ip)

		resp, err := http.Get("http://"+ip+"/cat/"+hash)
		if err != nil {
			log.Print(err)
		}

		io.Copy(w, resp.Body)
		defer resp.Body.Close()
	} else {
		// Writer the body to file
		_, err = io.Copy(w, out)
		if err != nil  {
			log.Fatal(err)
		}
	}

}

func Store(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		path := r.FormValue("path")
		filename := filepath.Base(path)
		filename = strings.ToLower(filename)
		hash := HashStr(filename)
		newPurgeInfo := d7024e.PurgeInformation{Key: hash, Pinned: false, TimeToLive:d7024e.TTL}
		kademlia.PublishData(newPurgeInfo, path)
		log.Println("handlers.go Store(): path = " + path)
		fmt.Fprint(w, hash)
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}

func Pin(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		vars := mux.Vars(r)
		//kademlia.Pin(hash)
		fmt.Fprintln(w, vars["hash"])
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}

func Unpin(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		vars := mux.Vars(r)
		//kademlia.Unpin(hash)
		fmt.Fprintln(w, vars["hash"])
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}

func AddNode(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		vars := mux.Vars(r)
		addr := vars["addr"]
		bootstrap := r.FormValue("bootstrap")

		kademlia := d7024e.NewKademlia()
		go kademlia.Run(bootstrap, addr)
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	fp := kademlia.Get(hash)
	log.Println(fp)
	filename := filepath.Base(fp)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	// Create the file
	out, err := os.Open(fp)
	if err != nil  {
		log.Fatal(err)
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(w, out)
	if err != nil  {
		log.Fatal(err)
	}
}

func isLocalHost(r *http.Request) bool{
	u, err := url.Parse("http://"+r.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	switch u.Hostname(){
	case "127.0.0.1":
		fallthrough
	case "::1":
		return true
	default:
		return false
	}
}
