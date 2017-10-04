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

	"github.com/gorilla/mux"
	"Mobile_D7024E/d7024e"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kademlia API\n\nLocalhost:\n     GET: /cat/{hash}\n     POST: /store/\n     GET: /pin/{hash}\n     GET: /unpin/{hash}\n     GET:/addnode/{addr}?boostrap={bootstrap_addr}\nPublic:\n     GET: /download/{hash}")
}

func Cat(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		vars := mux.Vars(r)
		filename := vars["filename"]
		hash := HashStr(filename)
		path := kademlia.Get(hash)
		fmt.Fprintln(w, path)
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
	
}

func Store(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		path := r.FormValue("path")
		filename := filepath.Base(path)
		filename = strings.ToLower(filename)
		hash := HashStr(filename)
		kademlia.Set(hash, path)
		//kademlia.PublishData(hash, path)
		fmt.Fprint(w, hash)
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}	

func Pin(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		fmt.Fprintln(w, "Pin endpoint")
		vars := mux.Vars(r)
		fmt.Fprintln(w, vars["hash"])
	} else{
		fmt.Fprintln(w, "Localhost only.")
	}
}

func Unpin(w http.ResponseWriter, r *http.Request) {
	if isLocalHost(r){
		fmt.Fprintln(w, "Unpin endpoint")
		vars := mux.Vars(r)
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

	// Get the data
	/*resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()*/

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
