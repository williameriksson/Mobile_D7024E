package api

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"strconv"
	"path/filepath"

	"github.com/gorilla/mux"
	"Mobile_D7024E/d7024e"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kademlia API:\nGET: /cat/{hash}\nPOST: /store/\nGET: /pin/{hash}\nGET: /unpin/{hash}\n GET:/addnode/{addr}?boostrap={addr}")
}

func Cat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	hash := HashStr(filename)
	path := kademlia.Get(hash)
	fmt.Fprintln(w, path)
}

func Store(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	filename := filepath.Base(path)
	hash := HashStr(filename)
	kademlia.Store(hash, path)
	fmt.Fprint(w, hash)
}	

func Pin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pin endpoint")
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars["hash"])
}

func Unpin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Unpin endpoint")
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars["hash"])
}

func AddNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["addr"]
	bootstrap := r.FormValue("bootstrap")

	kademlia := d7024e.NewKademlia()
	go kademlia.Run(bootstrap, addr)
}
